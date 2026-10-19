// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// Where a resolved dataplane writer user ID came from.
const (
	WriterSourceConfig       = "config"
	WriterSourceServiceToken = "service token"
)

// ServiceUserID returns the Keystone user ID that the provider client's token
// was issued for.
func ServiceUserID(client *gophercloud.ProviderClient) (string, error) {
	if client == nil {
		return "", errors.New("no provider client")
	}
	// Both tokens.CreateResult and tokens.GetResult implement this.
	result, ok := client.GetAuthResult().(interface {
		ExtractUser() (*tokens.User, error)
	})
	if !ok {
		return "", fmt.Errorf("unexpected auth result of type %T, cannot read the token's user", client.GetAuthResult())
	}
	user, err := result.ExtractUser()
	if err != nil {
		return "", fmt.Errorf("cannot read the user from the service token: %w", err)
	}
	if user == nil || user.ID == "" {
		return "", errors.New("the service token has no user ID")
	}
	return user.ID, nil
}

// ResolveWriterUserID picks the dataplane writer user ID: the configured value
// if there is one, else the service user's ID. It also returns where the ID
// came from. The ID is empty when neither is available.
func ResolveWriterUserID(configured, serviceUserID string) (id, source string) {
	if configured != "" {
		return configured, WriterSourceConfig
	}
	return serviceUserID, WriterSourceServiceToken
}
