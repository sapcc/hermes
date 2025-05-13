// Copyright 2022 SAP SE
// SPDX-FileCopyrightText: 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sapcc/go-bits/gopherpolicy"

	"github.com/sapcc/hermes/pkg/storage"
)

type v1Provider struct {
	validator   gopherpolicy.Validator
	storage     storage.Storage
	versionData VersionData
}

// NewV1Handler creates a http.Handler that serves the Hermes v1 API.
// It also returns the VersionData for this API version which is needed for the
// version advertisement on "GET /".
func NewV1Handler(validator gopherpolicy.Validator, storageInterface storage.Storage) (http.Handler, VersionData) {
	r := mux.NewRouter()

	p := &v1Provider{
		validator: validator,
		storage:   storageInterface,
	}
	p.versionData = VersionData{
		Status: "CURRENT",
		ID:     "v1",
		Links: []versionLinkData{
			{
				Relation: "self",
				URL:      p.Path(),
			},
			{
				Relation: "describedby",
				URL:      "https://github.com/sapcc/hermes/tree/master/docs",
				Type:     "text/html",
			},
		},
	}

	r.Methods("GET").Path("/v1/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		p.versionData.Links[0].URL = fmt.Sprintf("%s://%s%s/", getProtocol(req), req.Host, p.Path())
		ReturnJSON(res, 200, map[string]any{"version": p.versionData})
	})

	r.Methods("GET").Path("/v1/events").HandlerFunc(
		observeDuration(observeResponseSize(p.ListEvents, "ListEvents"), "ListEvents"))
	r.Methods("GET").Path("/v1/events/{event_id}").HandlerFunc(
		observeDuration(observeResponseSize(p.GetEventDetails, "GetEventDetails"), "GetEventDetails"))
	r.Methods("GET").Path("/v1/attributes/{attribute_name}").HandlerFunc(
		observeDuration(observeResponseSize(p.GetAttributes, "GetAttributes"), "GetAttributes"))

	return r, p.versionData
}

// Path constructs a full URL for a given URL path below the /v1/ endpoint.
func (p *v1Provider) Path(elements ...string) string {
	parts := []string{
		strings.TrimSuffix( /*p.Driver.Cluster().Config.CatalogURL*/ "", "/"),
		"v1",
	}
	parts = append(parts, elements...)
	return strings.Join(parts, "/")
}

// AuthHandler wraps endpoint handlers with consistent auth logic.
func (p *v1Provider) AuthHandler(w http.ResponseWriter, r *http.Request, rule string) (*gopherpolicy.Token, bool) {
	token := p.validator.CheckToken(r)

	// Initialize request context with URL vars
	token.Context.Request = mux.Vars(r)

	// Handle domain_id with form value priority
	if formDomainID := r.FormValue("domain_id"); formDomainID != "" {
		token.Context.Request["domain_id"] = formDomainID
	} else {
		token.Context.Request["domain_id"] = token.Context.Auth["domain_id"]
	}

	// Handle project_id with form value priority
	if formProjectID := r.FormValue("project_id"); formProjectID != "" {
		token.Context.Request["project_id"] = formProjectID
	} else {
		token.Context.Request["project_id"] = token.Context.Auth["project_id"]
	}

	ok := token.Require(w, rule)
	return token, ok
}
