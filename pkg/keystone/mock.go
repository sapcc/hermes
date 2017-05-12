package keystone

import (
	"github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud"
)

type mock struct{}

func Mock() Interface {
	return mock{}
}

func (d mock) keystoneClient() (*gophercloud.ServiceClient, error) {
	return nil, nil
}

func (d mock) Client() *gophercloud.ProviderClient {
	return nil
}

//ListDomains implements the Driver interface.
func (d mock) ListDomains() ([]KeystoneDomain, error) {
	return nil, nil
}

//ListProjects implements the Driver interface.
func (d mock) ListProjects() ([]KeystoneProject, error) {
	return nil, nil
}

//CheckUserPermission implements the Driver interface.
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
	return "ccdemo", nil
}

func (d mock) UserName(id string) (string, error) {
	return "I056593", nil
}
