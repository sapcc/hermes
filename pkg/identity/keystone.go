/*******************************************************************************
*
* Copyright 2022 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package identity

import (
	"context"
	"fmt"

	policy "github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/spf13/viper"

	"github.com/sapcc/go-bits/logg"
)

// Keystone Openstack Keystone implementation
type Keystone struct{}

// The JSON mappings here are for parsing Keystone responses
type keystoneToken struct {
	DomainScope  keystoneTokenThing         `json:"domain"`
	ProjectScope keystoneTokenThingInDomain `json:"project"`
	Roles        []keystoneTokenThing       `json:"roles"`
	User         keystoneTokenThingInDomain `json:"user"`
	ExpiresAt    string                     `json:"expires_at"`
}

// The JSON mappings here are for parsing Keystone responses
type keystoneTokenThing struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// The JSON mappings here are for parsing Keystone responses
type keystoneTokenThingInDomain struct {
	keystoneTokenThing
	Domain keystoneTokenThing `json:"domain"`
}

func (d Keystone) keystoneClient(ctx context.Context) (*gophercloud.ServiceClient, error) {
	logg.Debug("Getting service user Identity token...")
	if providerClient == nil {
		var err error
		// providerClient, err = openstack.NewClient(viper.GetString("Keystone.auth_url"))
		opts := d.AuthOptions()
		providerClient, err = openstack.AuthenticatedClient(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("cannot initialize OpenStack client: %w", err)
		}
	}

	//TODO: crashes with RegionName != ""
	return openstack.NewIdentityV3(providerClient,
		gophercloud.EndpointOpts{Region: "", Availability: gophercloud.AvailabilityPublic},
	)
}

// Client for Keystone connection
func (d Keystone) Client() (*gophercloud.ProviderClient, error) {
	var kc Keystone

	err := viper.UnmarshalKey("Keystone", &kc)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return nil, nil
}

// ValidateToken checks a token with Keystone
func (d Keystone) ValidateToken(ctx context.Context, token string) (policy.Context, error) {
	cachedToken := getCachedToken(tokenCache, token)
	if cachedToken != nil {
		return cachedToken.ToContext(), nil
	}

	client, err := d.keystoneClient(ctx)
	if err != nil {
		return policy.Context{}, err
	}

	response := tokens.Get(ctx, client, token)
	if response.Err != nil {
		// this includes 4xx responses, so after this point, we can be sure that the token is valid
		return policy.Context{}, response.Err
	}

	// use a custom token struct instead of tMap.Token which is way incomplete
	var tokenData keystoneToken
	err = response.ExtractInto(&tokenData)
	if err != nil {
		return policy.Context{}, err
	}
	d.updateCaches(&tokenData, token)
	return tokenData.ToContext(), nil
}

// updateCaches fills caches for Keystone lookups
func (d Keystone) updateCaches(token *keystoneToken, tokenStr string) {
	addTokenToCache(tokenCache, tokenStr, token)
}

// ToContext
func (t *keystoneToken) ToContext() policy.Context {
	c := policy.Context{
		Roles: make([]string, 0, len(t.Roles)),
		Auth: map[string]string{
			"user_id":             t.User.ID,
			"user_name":           t.User.Name,
			"user_domain_id":      t.User.Domain.ID,
			"user_domain_name":    t.User.Domain.Name,
			"domain_id":           t.DomainScope.ID,
			"domain_name":         t.DomainScope.Name,
			"project_id":          t.ProjectScope.ID,
			"project_name":        t.ProjectScope.Name,
			"project_domain_id":   t.ProjectScope.Domain.ID,
			"project_domain_name": t.ProjectScope.Domain.Name,
			"tenant_id":           t.ProjectScope.ID,
			"tenant_name":         t.ProjectScope.Name,
			"tenant_domain_id":    t.ProjectScope.Domain.ID,
			"tenant_domain_name":  t.ProjectScope.Domain.Name,
		},
		Request: nil,
		Logger:  logg.Debug,
	}
	for key, value := range c.Auth {
		if value == "" {
			delete(c.Auth, key)
		}
	}
	for _, role := range t.Roles {
		c.Roles = append(c.Roles, role.Name)
	}
	if c.Request == nil {
		c.Request = map[string]string{}
	}

	return c
}

// AuthOptions fills in Keystone options with hermes config values
func (d Keystone) AuthOptions() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: viper.GetString("Keystone.auth_url"),
		Username:         viper.GetString("Keystone.username"),
		Password:         viper.GetString("Keystone.password"),
		DomainName:       viper.GetString("Keystone.user_domain_name"),
		Scope: &gophercloud.AuthScope{
			ProjectName: viper.GetString("Keystone.project_name"),
			DomainName:  viper.GetString("Keystone.project_domain_name"),
		},
		AllowReauth: true,
	}
}
