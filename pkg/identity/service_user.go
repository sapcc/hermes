// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// Where a resolved dataplane writer user ID came from.
const (
	WriterSourceConfig       = "config"
	WriterSourceServiceToken = "service token"
)

// ServiceUser identifies the Keystone user of Hermes' own service token.
type ServiceUser struct {
	ID     string
	Name   string
	Domain string
}

// String formats the user for log lines, e.g. "hermes@Default".
func (u ServiceUser) String() string {
	if u.Name == "" {
		return "unknown name"
	}
	if u.Domain == "" {
		return u.Name
	}
	return u.Name + "@" + u.Domain
}

// ReadServiceUser returns the Keystone user that the provider client's token
// was issued for.
func ReadServiceUser(client *gophercloud.ProviderClient) (ServiceUser, error) {
	if client == nil {
		return ServiceUser{}, errors.New("no provider client")
	}
	// Both tokens.CreateResult and tokens.GetResult implement this.
	result, ok := client.GetAuthResult().(interface {
		ExtractUser() (*tokens.User, error)
	})
	if !ok {
		return ServiceUser{}, fmt.Errorf("unexpected auth result of type %T, cannot read the token's user", client.GetAuthResult())
	}
	user, err := result.ExtractUser()
	if err != nil {
		return ServiceUser{}, fmt.Errorf("cannot read the user from the service token: %w", err)
	}
	if user == nil || user.ID == "" {
		return ServiceUser{}, errors.New("the service token has no user ID")
	}
	domain := user.Domain.Name
	if domain == "" {
		domain = user.Domain.ID
	}
	return ServiceUser{ID: user.ID, Name: user.Name, Domain: domain}, nil
}

// ResolveWriterUserID picks the dataplane writer user ID: the configured value
// if it is not blank, else the service user's ID. A configured value is
// trimmed. It also returns where the ID came from. The ID is empty when
// neither is available.
func ResolveWriterUserID(configured, serviceUserID string) (id, source string) {
	if c := strings.TrimSpace(configured); c != "" {
		return c, WriterSourceConfig
	}
	return serviceUserID, WriterSourceServiceToken
}
