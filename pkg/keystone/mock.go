package keystone

import (
	"github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud"
	"github.com/spf13/viper"
)

type mock struct{}

// Mock keystone implementation
func Mock() Driver {
	return mock{}
}

func (d mock) keystoneClient() (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d mock) Client() *gophercloud.ProviderClient {
	return nil
}

func (d mock) ValidateToken(token string) (policy.Context, error) {

	return policy.Context{}, nil
}

func (d mock) Authenticate(credentials *gophercloud.AuthOptions) (policy.Context, error) {
	return policy.Context{}, nil
}

func (d mock) DomainName(id string) (string, error) {
	return "monsoon3", nil
}

func (d mock) ProjectName(id string) (string, error) {
	return "ceilometer-cadf-delete-me", nil
}

func (d mock) UserName(id string) (string, error) {
	return "I056593", nil
}

func (d mock) UserId(name string) (string, error) {
	return "eb5cd8f904b06e8b2a6eb86c8b04c08e6efb89b92da77905cc8c475f30b0b812", nil
}

func (d mock) RoleName(id string) (string, error) {
	return "audit_viewer", nil
}

func (d mock) GroupName(id string) (string, error) {
	return "admins", nil
}

func (d mock) AuthOptions() *gophercloud.AuthOptions {
	return &gophercloud.AuthOptions{
		IdentityEndpoint: viper.GetString("keystone.auth_url"),
		Username:         viper.GetString("keystone.username"),
		Password:         viper.GetString("keystone.password"),
		DomainName:       viper.GetString("keystone.user_domain_name"),
		// Note: gophercloud only allows for user & project in the same domain
		TenantName: viper.GetString("keystone.project_name"),
	}
}
