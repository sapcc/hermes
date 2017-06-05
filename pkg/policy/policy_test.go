package policy

import (
	"testing"

	policy "github.com/databus23/goslo.policy"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/stretchr/testify/assert"
)


// Setup test
func GetEnforcer() *policy.Enforcer {
	const path = "../../etc/policy.json"


	policyenforcer, err := util.LoadPolicyFile(path)

	if err != nil {
		return nil
	}
	return policyenforcer
}

func Test_Policy_AuditViewerTrue(t *testing.T) {
	policyenforcer := GetEnforcer()
	c := policy.Context{
		Roles: []string{
			"audit_viewer",
		},
		Auth: map[string]string{
			"user_id":             "aaaa",
			"user_name":           "aaaa",
			"user_domain_id":      "aaaa",
			"user_domain_name":    "aaaa",
			"domain_id":           "ca1b267e149d4e44bf53d28d1c8d6bc9",
			"domain_name":         "aaaa",
			"project_id":          "7a09c05926ec452ca7992af4aa03c31d",
			"project_name":        "aaaa",
			"project_domain_id":   "aaaa",
			"project_domain_name": "aaaa",
			"tenant_id":           "aaaa",
			"tenant_name":         "aaaa",
			"tenant_domain_id":    "aaaa",
			"tenant_domain_name":  "aaaa",
		},
		Request: map[string]string{
			"domain_id":           "ca1b267e149d4e44bf53d28d1c8d6bc9",
			"project_id":          "7a09c05926ec452ca7992af4aa03c31d",
		},
		Logger:  util.LogDebug,
	}
	assert.True(t, policyenforcer.Enforce("event:show", c))
}

func Test_Policy_UnknownRoleFalse(t *testing.T) {
	policyenforcer := GetEnforcer()
	c := policy.Context{
		Roles: []string{
			"unknown_role",
		},
		Auth: map[string]string{
			"user_id":             "aaaa",
			"user_name":           "aaaa",
			"user_domain_id":      "aaaa",
			"user_domain_name":    "aaaa",
			"domain_id":           "ca1b267e149d4e44bf53d28d1c8d6bc9",
			"domain_name":         "aaaa",
			"project_id":          "7a09c05926ec452ca7992af4aa03c31d",
			"project_name":        "aaaa",
			"project_domain_id":   "aaaa",
			"project_domain_name": "aaaa",
			"tenant_id":           "aaaa",
			"tenant_name":         "aaaa",
			"tenant_domain_id":    "aaaa",
			"tenant_domain_name":  "aaaa",
		},
		Request: map[string]string{
			"domain_id":           "ca1b267e149d4e44bf53d28d1c8d6bc9",
			"project_id":          "7a09c05926ec452ca7992af4aa03c31d",
		},
		Logger:  util.LogDebug,
	}
	assert.False(t, policyenforcer.Enforce("event:show", c))
}
