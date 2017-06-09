package identity

import (
	"github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud"
	"github.com/spf13/viper"
)

type Mock struct{}

func (d Mock) keystoneClient() (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d Mock) Client() *gophercloud.ProviderClient {
	return nil
}

func (d Mock) ValidateToken(token string) (policy.Context, error) {

	return policy.Context{}, nil
}

func (d Mock) Authenticate(credentials *gophercloud.AuthOptions) (policy.Context, error) {
	return policy.Context{}, nil
}

func (d Mock) DomainName(id string) (string, error) {
	return "monsoon3", nil
}

func (d Mock) ProjectName(id string) (string, error) {
	return "ceilometer-cadf-delete-me", nil
}

func (d Mock) UserName(id string) (string, error) {
	return "I056593", nil
}

func (d Mock) UserId(name string) (string, error) {
	return "eb5cd8f904b06e8b2a6eb86c8b04c08e6efb89b92da77905cc8c475f30b0b812", nil
}

func (d Mock) RoleName(id string) (string, error) {
	return "audit_viewer", nil
}

func (d Mock) GroupName(id string) (string, error) {
	return "admins", nil
}

func (d Mock) AuthOptions() *gophercloud.AuthOptions {
	return &gophercloud.AuthOptions{
		IdentityEndpoint: viper.GetString("Keystone.auth_url"),
		Username:         viper.GetString("Keystone.username"),
		Password:         viper.GetString("Keystone.password"),
		DomainName:       viper.GetString("Keystone.user_domain_name"),
		// Note: gophercloud only allows for user & project in the same domain
		TenantName: viper.GetString("Keystone.project_name"),
	}
}
