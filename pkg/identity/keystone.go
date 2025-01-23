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

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/spf13/viper"

	"github.com/sapcc/go-bits/gopherpolicy"
)

// NewTokenValidator connects to Keystone using the provided OpenStack
// credentials and constructs a gopherpolicy.TokenValidator instance.
func NewTokenValidator(ctx context.Context) (*gopherpolicy.TokenValidator, error) {
	opts := authOptions()
	providerClient, err := openstack.AuthenticatedClient(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize OpenStack client: %w", err)
	}

	//TODO: crashes with RegionName != ""
	identityV3, err := openstack.NewIdentityV3(providerClient,
		gophercloud.EndpointOpts{Region: "", Availability: gophercloud.AvailabilityPublic},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize Keystone client: %w", err)
	}

	tv := gopherpolicy.TokenValidator{
		IdentityV3: identityV3,
		Cacher:     gopherpolicy.InMemoryCacher(),
	}
	err = tv.LoadPolicyFile(viper.GetString("hermes.PolicyFilePath"), nil)
	if err != nil {
		return nil, err
	}

	return &tv, nil
}

func authOptions() gophercloud.AuthOptions {
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
