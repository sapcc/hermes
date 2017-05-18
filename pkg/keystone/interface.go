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
	policy "github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud"
)

// Driver is an interface that wraps the authentication of the service user and
// token checking of API users. Because it is an interface, the real implementation
// can be mocked away in unit tests.
type Interface interface {
	//Return the main gophercloud client from which the respective service
	//clients can be derived. For mock drivers, this returns nil, so test code
	//should be prepared to handle a nil Client() where appropriate.
	Client() *gophercloud.ProviderClient
	AuthOptions() *gophercloud.AuthOptions
	/********** requests to Keystone **********/
	ListDomains() ([]KeystoneDomain, error)
	ListProjects() ([]KeystoneProject, error)
	ValidateToken(token string) (policy.Context, error)
	Authenticate(credentials *gophercloud.AuthOptions) (policy.Context, error)
	DomainName(id string) (string, error)
	ProjectName(id string) (string, error)
	UserName(id string) (string, error)
}

//KeystoneDomain describes the basic attributes of a Keystone domain.
type KeystoneDomain struct {
	UUID string `json:"id"`
	Name string `json:"name"`
}

//KeystoneProject describes the basic attributes of a Keystone project.
type KeystoneProject struct {
	UUID string `json:"id"`
	Name string `json:"name"`
}

//KeystoneProject describes the basic attributes of a Keystone project.
type KeystoneUser struct {
	UUID string `json:"id"`
	Name string `json:"name"`
}
