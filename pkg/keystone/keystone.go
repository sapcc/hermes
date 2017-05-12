/*******************************************************************************
*
* Copyright 2017 SAP SE
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

package keystone

import (
	"fmt"

	policy "github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
)

func Keystone() Interface {
	return keystone{}
}

type keystone struct {
	ProviderClient *gophercloud.ProviderClient
}

func (d keystone) keystoneClient() (*gophercloud.ServiceClient, error) {
	if d.ProviderClient == nil {
		var err error
		d.ProviderClient, err = openstack.NewClient(viper.GetString("keystone.auth_url"))
		if err != nil {
			return nil, fmt.Errorf("cannot initialize OpenStack client: %v", err)
		}
	}

	return openstack.NewIdentityV3(d.ProviderClient,
		gophercloud.EndpointOpts{Availability: gophercloud.AvailabilityPublic},
	)
}

func (d keystone) Client() *gophercloud.ProviderClient {
	var kc keystone

	err := viper.UnmarshalKey("keystone", &kc)
	if err != nil {
		fmt.Println("unable to decode into struct, %v", err)
	}

	return nil
}

//ListDomains implements the Driver interface.
func (d keystone) ListDomains() ([]KeystoneDomain, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return nil, err
	}

	//gophercloud does not support domain listing yet - do it manually
	url := client.ServiceURL("domains")
	var result gophercloud.Result
	_, err = client.Get(url, &result.Body, nil)
	if err != nil {
		return nil, err
	}

	var data struct {
		Domains []KeystoneDomain `json:"domains"`
	}
	err = result.ExtractInto(&data)
	return data.Domains, err
}

//ListProjects implements the Driver interface.
func (d keystone) ListProjects() ([]KeystoneProject, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return nil, err
	}

	var result gophercloud.Result
	_, err = client.Get("/v3/auth/projects", &result.Body, nil)
	if err != nil {
		return nil, err
	}

	var data struct {
		Projects []KeystoneProject `json:"projects"`
	}
	err = result.ExtractInto(&data)
	return data.Projects, err
}

//CheckUserPermission implements the Driver interface.
func (d keystone) ValidateToken(token string) (policy.Context, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return policy.Context{}, err
	}

	response := tokens.Get(client, token)
	if response.Err != nil {
		//this includes 4xx responses, so after this point, we can be sure that the token is valid
		return policy.Context{}, response.Err
	}

	//use a custom token struct instead of tokens.Token which is way incomplete
	var tokenData keystoneToken
	err = response.ExtractInto(&tokenData)
	if err != nil {
		return policy.Context{}, err
	}
	return tokenData.ToContext(), nil
}

func (d keystone) Authenticate(credentials *gophercloud.AuthOptions) (policy.Context, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return policy.Context{}, err
	}
	response := tokens.Create(client, credentials)
	if response.Err != nil {
		//this includes 4xx responses, so after this point, we can be sure that the token is valid
		return policy.Context{}, response.Err
	}
	//use a custom token struct instead of tokens.Token which is way incomplete
	var tokenData keystoneToken
	err = response.ExtractInto(&tokenData)
	if err != nil {
		return policy.Context{}, err
	}
	return tokenData.ToContext(), nil
}

func (d keystone) DomainName(id string) (string, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	_, err = client.Get(fmt.Sprintf("/v3/domains/%s", id), &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Domain KeystoneDomain `json:"domain"`
	}
	err = result.ExtractInto(&data)
	return data.Domain.Name, err
}

func (d keystone) ProjectName(id string) (string, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	_, err = client.Get(fmt.Sprintf("/v3/projects/%s", id), &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		Project KeystoneProject `json:"project"`
	}
	err = result.ExtractInto(&data)
	return data.Project.Name, err
}

func (d keystone) UserName(id string) (string, error) {
	client, err := d.keystoneClient()
	if err != nil {
		return "", err
	}

	var result gophercloud.Result
	_, err = client.Get(fmt.Sprintf("/v3/users/%s", id), &result.Body, nil)
	if err != nil {
		return "", err
	}

	var data struct {
		User KeystoneUser `json:"user"`
	}
	err = result.ExtractInto(&data)
	return data.User.Name, err
}

type keystoneToken struct {
	DomainScope  keystoneTokenThing         `json:"domain"`
	ProjectScope keystoneTokenThingInDomain `json:"project"`
	Roles        []keystoneTokenThing       `json:"roles"`
	User         keystoneTokenThingInDomain `json:"user"`
}

type keystoneTokenThing struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type keystoneTokenThingInDomain struct {
	keystoneTokenThing
	Domain keystoneTokenThing `json:"domain"`
}

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
		Logger:  util.LogDebug,
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
