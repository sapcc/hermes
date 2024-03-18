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
//nolint:dupl
package identity

import (
	"fmt"

	"sync"

	policy "github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/sapcc/go-bits/logg"
)

// Keystone Openstack Keystone implementation
type Keystone struct {
	TokenRenewalMutex *sync.Mutex // Used for controlling the token refresh process
}

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

// keystoneNameID describes just the name and id of a Identity object.
//
//	The JSON mappings here are for parsing Keystone responses
type keystoneNameID struct {
	UUID string `json:"id"`
	Name string `json:"name"`
}

func (d Keystone) keystoneClient() (*gophercloud.ServiceClient, error) {
	logg.Debug("Getting service user Identity token...")
	if d.TokenRenewalMutex == nil {
		d.TokenRenewalMutex = &sync.Mutex{}
	}
	if providerClient == nil {
		var err error
		//providerClient, err = openstack.NewClient(viper.GetString("Keystone.auth_url"))
		opts := d.AuthOptions()
		providerClient, err = openstack.AuthenticatedClient(*opts)
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
func (d Keystone) ValidateToken(token string) (policy.Context, error) {
	cachedToken := getCachedToken(tokenCache, token)
	if cachedToken != nil {
		return cachedToken.ToContext(), nil
	}

	client, err := d.keystoneClient()
	if err != nil {
		return policy.Context{}, err
	}

	response := tokens.Get(client, token)
	if response.Err != nil {
		//this includes 4xx responses, so after this point, we can be sure that the token is valid
		return policy.Context{}, response.Err
	}

	//use a custom token struct instead of tMap.Token which is way incomplete
	var tokenData keystoneToken
	err = response.ExtractInto(&tokenData)
	if err != nil {
		return policy.Context{}, err
	}
	d.updateCaches(&tokenData, token)
	return tokenData.ToContext(), nil
}

// Authenticate with Keystone
func (d Keystone) Authenticate(credentials *gophercloud.AuthOptions) (policy.Context, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return policy.Context{}, err
	}
	response := tokens.Create(client, credentials)
	if response.Err != nil {
		//this includes 4xx responses, so after this point, we can be sure that the token is valid
		return policy.Context{}, response.Err
	}
	//use a custom token struct instead of tMap.Token which is way incomplete
	var tokenData keystoneToken
	err = response.ExtractInto(&tokenData)
	if err != nil {
		return policy.Context{}, err
	}
	return tokenData.ToContext(), nil
}

// DomainName with caching
func (d Keystone) DomainName(id string) (string, error) {
	cachedName, hit := getFromCache(domainNameCache, id)
	if hit {
		return cachedName, nil
	}

	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	url := client.ServiceURL("domains/" + id)
	_, err = client.Get(url, &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Domain keystoneNameID `json:"domain"`
	}
	err = result.ExtractInto(&data)
	if err == nil {
		updateCache(domainNameCache, id, data.Domain.Name)
	}
	return data.Domain.Name, err
}

// ProjectName with caching
func (d Keystone) ProjectName(id string) (string, error) {
	cachedName, hit := getFromCache(projectNameCache, id)
	if hit {
		return cachedName, nil
	}

	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	url := client.ServiceURL("projects/" + id)
	_, err = client.Get(url, &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Project keystoneNameID `json:"project"`
	}
	err = result.ExtractInto(&data)
	if err == nil {
		updateCache(projectNameCache, id, data.Project.Name)
	}
	return data.Project.Name, err
}

// UserName with Caching
func (d Keystone) UserName(id string) (string, error) {
	cachedName, hit := getFromCache(userNameCache, id)
	if hit {
		return cachedName, nil
	}

	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	url := client.ServiceURL("users/" + id)
	_, err = client.Get(url, &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		User keystoneNameID `json:"user"`
	}
	err = result.ExtractInto(&data)
	if err == nil {
		updateCache(userNameCache, id, data.User.Name)
		updateCache(userIDCache, data.User.Name, id)
	}
	return data.User.Name, err
}

// UserID with caching
func (d Keystone) UserID(name string) (string, error) {
	cachedID, hit := getFromCache(userIDCache, name)
	if hit {
		return cachedID, nil
	}

	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	url := client.ServiceURL("users?name=" + name)
	_, err = client.Get(url, &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		User []keystoneNameID `json:"user"`
	}
	err = result.ExtractInto(&data)
	userID := ""
	if err == nil {
		switch len(data.User) {
		case 0:
			err = errors.Errorf("No user found with name %s", name)
		case 1:
			userID = data.User[0].UUID
		default:
			logg.Info("Multiple users found with name %s - returning the first one", name)
			userID = data.User[0].UUID
		}
		updateCache(userIDCache, name, userID)
		updateCache(userNameCache, userID, name)
	}
	return userID, err
}

// RoleName with caching
func (d Keystone) RoleName(id string) (string, error) {
	cachedName, hit := getFromCache(roleNameCache, id)
	if hit {
		return cachedName, nil
	}

	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	url := client.ServiceURL("roles/" + id)
	_, err = client.Get(url, &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Role keystoneNameID `json:"role"`
	}
	err = result.ExtractInto(&data)
	if err == nil {
		updateCache(roleNameCache, id, data.Role.Name)
	}
	return data.Role.Name, err
}

// GroupName with caching
func (d Keystone) GroupName(id string) (string, error) {
	cachedName, hit := getFromCache(groupNameCache, id)
	if hit {
		return cachedName, nil
	}

	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	url := client.ServiceURL("groups/" + id)
	_, err = client.Get(url, &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Group keystoneNameID `json:"group"`
	}
	err = result.ExtractInto(&data)
	if err == nil {
		updateCache(groupNameCache, id, data.Group.Name)
	}
	return data.Group.Name, err
}

// updateCaches fills caches for Keystone lookups
func (d Keystone) updateCaches(token *keystoneToken, tokenStr string) {
	addTokenToCache(tokenCache, tokenStr, token)
	if token.DomainScope.ID != "" && token.DomainScope.Name != "" {
		updateCache(domainNameCache, token.DomainScope.ID, token.DomainScope.Name)
	}
	if token.ProjectScope.Domain.ID != "" && token.ProjectScope.Domain.Name != "" {
		updateCache(domainNameCache, token.ProjectScope.Domain.ID, token.ProjectScope.Domain.Name)
	}
	if token.ProjectScope.ID != "" && token.ProjectScope.Name != "" {
		updateCache(projectNameCache, token.ProjectScope.ID, token.ProjectScope.Name)
	}
	if token.User.ID != "" && token.User.Name != "" {
		updateCache(userNameCache, token.User.ID, token.User.Name)
		updateCache(userIDCache, token.User.Name, token.User.ID)
	}
	for _, role := range token.Roles {
		if role.ID != "" && role.Name != "" {
			updateCache(roleNameCache, role.ID, role.Name)
		}
	}
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

// RefreshToken fetches a new Identity auth token. It is also used
// to fetch the initial token on startup.
func (d Keystone) RefreshToken() error {
	//NOTE: This function is very similar to v3auth() in
	//gophercloud/openstack/client.go, but with a few differences:
	//
	//1. thread-safe token renewal
	//2. proper support for cross-domain scoping

	logg.Debug("Getting service user Identity token...")

	d.TokenRenewalMutex.Lock()
	defer d.TokenRenewalMutex.Unlock()

	providerClient.TokenID = ""

	//TODO: crashes with RegionName != ""
	eo := gophercloud.EndpointOpts{Region: ""}
	keystone, err := openstack.NewIdentityV3(providerClient, eo)
	if err != nil {
		return fmt.Errorf("cannot initialize Identity client: %w", err)
	}

	logg.Debug("Identity URL: %s", keystone.Endpoint)

	result := tokens.Create(keystone, d.AuthOptions())
	token, err := result.ExtractToken()
	if err != nil {
		return fmt.Errorf("cannot read token: %w", err)
	}
	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return fmt.Errorf("cannot read service catalog: %w", err)
	}

	providerClient.TokenID = token.ID
	providerClient.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
		return openstack.V3EndpointURL(catalog, opts)
	}

	return nil
}

// AuthOptions fills in Keystone options with hermes config values
func (d Keystone) AuthOptions() *gophercloud.AuthOptions {
	return &gophercloud.AuthOptions{
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
