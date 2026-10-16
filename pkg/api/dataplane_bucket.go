// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"

	"github.com/sapcc/go-bits/gopherpolicy"
	"github.com/sapcc/go-bits/logg"
)

const (
	// defaultDataplaneServiceType is the service catalog type of the Ceph RGW
	// Swift endpoint, used when hermes.dataplane_swift_service_type is not set.
	defaultDataplaneServiceType = "object-store-ceph"
	// dataplaneProbeObject is written and deleted on enable to check that the
	// project has object storage quota left.
	dataplaneProbeObject = ".hermes-enable-probe"
	// dataplaneBucketTimeout bounds all object store requests of one enable.
	dataplaneBucketTimeout = 30 * time.Second
)

// dataplaneBucketConfig is the operator configuration for bucket provisioning.
type dataplaneBucketConfig struct {
	// WriterUserID is the Keystone user ID of the log-router writer. It gets
	// read and write access to the customer's bucket through the container ACLs.
	WriterUserID string
	// ServiceType is the catalog type of the object store endpoint.
	ServiceType string
}

// bucketError is a failure of ensureDataplaneBucket that maps to an API response.
type bucketError struct {
	Status  int
	Message string
}

func (e *bucketError) Error() string { return e.Message }

// objectStoreOp names the step of the provisioning that got an error status.
type objectStoreOp int

const (
	opHeadContainer objectStoreOp = iota
	opCreateContainer
	opSetACL
	opProbe
)

// ensureDataplaneBucket makes sure the container exists in the token's own
// project, carries the ACLs that let log-router read and write it, and that the
// project has quota left for log-router's objects.
//
// All requests use the caller's own token: the project that creates the
// container owns it (quota and billing follow the owner), and a container POST
// that sets ACLs makes the requester the owner, so it must be the same token.
func (p *v1Provider) ensureDataplaneBucket(ctx context.Context, token *gopherpolicy.Token, projectID, bucket string) *bucketError {
	if projectID == "" || token.Context.Auth["project_id"] != projectID {
		return &bucketError{http.StatusConflict, "enable dataplane with a token scoped to the project; the bucket must be created by the project"}
	}
	if p.dataplaneBucket.WriterUserID == "" {
		logg.Error("dataplane-config PUT: hermes.dataplane_writer_user_id (HERMES_DATAPLANE_WRITER_USER_ID) is not set, cannot enable dataplane for project %s", projectID)
		return &bucketError{http.StatusServiceUnavailable, "dataplane bucket provisioning is not configured on this Hermes instance"}
	}
	if token.ProviderClient == nil || token.ProviderClient.EndpointLocator == nil {
		logg.Error("dataplane-config PUT: token for project %s has no service catalog access", projectID)
		return &bucketError{http.StatusBadGateway, "cannot determine the object storage endpoint"}
	}
	serviceType := p.dataplaneBucket.ServiceType
	if serviceType == "" {
		serviceType = defaultDataplaneServiceType
	}
	endpoint, err := token.ProviderClient.EndpointLocator(gophercloud.EndpointOpts{
		Type:         serviceType,
		Availability: gophercloud.AvailabilityPublic,
	})
	if notFound, ok := errors.AsType[*gophercloud.ErrEndpointNotFound](err); ok {
		logg.Error("dataplane-config PUT: no public %q endpoint for project %s: %s", serviceType, projectID, notFound)
		return &bucketError{http.StatusServiceUnavailable, "object storage for dataplane is not available in this region"}
	}
	if err != nil {
		logg.Error("dataplane-config PUT: cannot look up the %q endpoint for project %s: %s", serviceType, projectID, err)
		return &bucketError{http.StatusBadGateway, "cannot determine the object storage endpoint"}
	}
	// Defense in depth: the token is only ever sent to the caller's own account.
	account, err := url.Parse(endpoint)
	if err != nil || !strings.HasSuffix(strings.TrimSuffix(account.Path, "/"), "/AUTH_"+projectID) {
		logg.Error("dataplane-config PUT: refusing %q endpoint %q for project %s: not the project's own account", serviceType, endpoint, projectID)
		return &bucketError{http.StatusBadGateway, "cannot determine the object storage endpoint"}
	}

	ctx, cancel := context.WithTimeout(ctx, dataplaneBucketTimeout)
	defer cancel()

	client := p.bucketHTTPClient
	if client == nil {
		client = defaultBucketHTTPClient
	}
	store := objectStore{client: client, tokenID: token.ProviderClient.TokenID, projectID: projectID}
	containerURL := account.JoinPath(bucket).String()
	probeURL := account.JoinPath(bucket, dataplaneProbeObject).String()

	// existed is true when the container was there before this request. A 403
	// on changing it then means it belongs to another project.
	existed := false
	fail := func(op objectStoreOp, method, target string, reply *objectStoreReply) *bucketError {
		logg.Error("dataplane-config PUT: %s %s for project %s returned %d", method, target, projectID, reply.StatusCode)
		return mapObjectStoreStatus(op, reply.StatusCode, existed, bucket)
	}

	head, berr := store.do(ctx, http.MethodHead, containerURL, nil, nil)
	if berr != nil {
		return berr
	}
	var currentRead, currentWrite string
	switch {
	case head.StatusCode == http.StatusNotFound:
		put, putErr := store.do(ctx, http.MethodPut, containerURL, nil, nil)
		if putErr != nil {
			return putErr
		}
		if !isSuccess(put.StatusCode) {
			return fail(opCreateContainer, http.MethodPut, containerURL, put)
		}
	case isSuccess(head.StatusCode):
		existed = true
		currentRead = head.Header.Get("X-Container-Read")
		currentWrite = head.Header.Get("X-Container-Write")
	default:
		return fail(opHeadContainer, http.MethodHead, containerURL, head)
	}

	// The writer needs read access as well: log-router reads the previous
	// digest to chain the next one, and RGW answers 403 instead of 404 for a
	// missing object without it.
	writerGrant := "*:" + p.dataplaneBucket.WriterUserID
	wantRead := mergeACL(currentRead, projectID+":*", writerGrant)
	wantWrite := mergeACL(currentWrite, writerGrant)
	// Only POST when something is missing. Any ACL POST makes the requester the
	// owner, and on a container that still belongs to another project it fails,
	// even when the grants are already in place.
	if !sameACL(wantRead, currentRead) || !sameACL(wantWrite, currentWrite) {
		acl := http.Header{}
		acl.Set("X-Container-Read", wantRead)
		acl.Set("X-Container-Write", wantWrite)
		post, postErr := store.do(ctx, http.MethodPost, containerURL, acl, nil)
		if postErr != nil {
			return postErr
		}
		if !isSuccess(post.StatusCode) {
			return fail(opSetACL, http.MethodPost, containerURL, post)
		}
	}

	// Creating a container does not check byte quota, writing an object does.
	// Write a tiny object so a project without quota fails now, not later in
	// log-router where its events would be dropped.
	probeHeader := http.Header{}
	probeHeader.Set("Content-Type", "application/octet-stream")
	probe, berr := store.do(ctx, http.MethodPut, probeURL, probeHeader, []byte{'\n'})
	if berr != nil {
		return berr
	}
	if !isSuccess(probe.StatusCode) {
		return fail(opProbe, http.MethodPut, probeURL, probe)
	}
	del, berr := store.do(ctx, http.MethodDelete, probeURL, nil, nil)
	if berr != nil {
		return berr
	}
	if !isSuccess(del.StatusCode) && del.StatusCode != http.StatusNotFound {
		return fail(opProbe, http.MethodDelete, probeURL, del)
	}
	return nil
}

// objectStore sends requests to the object store with the caller's token.
type objectStore struct {
	client    *http.Client
	tokenID   string
	projectID string
}

// objectStoreReply is the part of an object store response that we look at.
type objectStoreReply struct {
	StatusCode int
	Header     http.Header
}

func (s objectStore) do(ctx context.Context, method, target string, header http.Header, body []byte) (*objectStoreReply, *bucketError) {
	var reqBody io.Reader = http.NoBody
	if body != nil {
		reqBody = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, target, reqBody)
	if err != nil {
		logg.Error("dataplane-config PUT: cannot build %s request for project %s: %s", method, s.projectID, err)
		return nil, &bucketError{http.StatusBadGateway, "object storage request failed"}
	}
	maps.Copy(req.Header, header)
	req.Header.Set("X-Auth-Token", s.tokenID)
	resp, err := s.client.Do(req)
	if err != nil {
		logg.Error("dataplane-config PUT: %s %s failed for project %s: %s", method, target, s.projectID, err)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &bucketError{http.StatusBadGateway, "object storage did not answer in time"}
		}
		return nil, &bucketError{http.StatusBadGateway, "object storage request failed"}
	}
	defer resp.Body.Close()
	return &objectStoreReply{resp.StatusCode, resp.Header.Clone()}, nil
}

func isSuccess(status int) bool { return status/100 == 2 }

// mapObjectStoreStatus converts a failing object store status into an API
// error. existed reports whether the container was there before the request.
func mapObjectStoreStatus(op objectStoreOp, status int, existed bool, bucket string) *bucketError {
	switch {
	case status == http.StatusUnauthorized:
		// RGW answers 401 for a token without an accepted role.
		return &bucketError{http.StatusForbidden, "your token needs the objectstore_admin role in this project"}
	case status == http.StatusForbidden && existed && (op == opSetACL || op == opProbe):
		return &bucketError{http.StatusConflict, fmt.Sprintf("container %s exists but is owned by another project; ask the operator to migrate it or choose another target_bucket", bucket)}
	case status == http.StatusForbidden:
		return &bucketError{http.StatusForbidden, "object storage denied the request"}
	case status == http.StatusBadRequest && op == opCreateContainer:
		// RGW answers 400 TooManyBuckets when the container limit is reached.
		return &bucketError{http.StatusConflict, "this project cannot create more containers"}
	case status == http.StatusConflict:
		return &bucketError{http.StatusConflict, "container name is not available"}
	case status == http.StatusRequestEntityTooLarge:
		return &bucketError{http.StatusConflict, "project has no free object storage quota"}
	default:
		return &bucketError{http.StatusBadGateway, fmt.Sprintf("object storage returned unexpected status %d", status)}
	}
}

// aclEntries splits a comma-separated ACL into its entries, trimmed and
// without empty entries or duplicates, in their original order.
func aclEntries(acl string) []string {
	var entries []string
	for e := range strings.SplitSeq(acl, ",") {
		if e = strings.TrimSpace(e); e != "" && !slices.Contains(entries, e) {
			entries = append(entries, e)
		}
	}
	return entries
}

// mergeACL returns the comma-separated ACL existing with the missing entries
// appended. Existing entries keep their order.
func mergeACL(existing string, add ...string) string {
	entries := aclEntries(existing)
	for _, entry := range add {
		if !slices.Contains(entries, entry) {
			entries = append(entries, entry)
		}
	}
	return strings.Join(entries, ",")
}

// sameACL reports whether two ACLs hold the same entries. RGW returns ACL
// entries in its own order, so the order does not count.
func sameACL(a, b string) bool {
	ea, eb := aclEntries(a), aclEntries(b)
	slices.Sort(ea)
	slices.Sort(eb)
	return slices.Equal(ea, eb)
}

// defaultBucketHTTPClient never follows redirects, so the caller's token is
// not forwarded to a host the catalog did not name. The overall deadline comes
// from the request context.
var defaultBucketHTTPClient = &http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}
