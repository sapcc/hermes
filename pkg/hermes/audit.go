package hermes

import (
	"github.com/sapcc/hermes/pkg/configdb"
	"github.com/sapcc/hermes/pkg/util"
)

type AuditDetail struct {
	Enabled  bool   `json:"enabled"`
	TenantID string `json:"TenantID"`
}

//GetAudit returns the config for auditing matching a tenant in JSON
func GetAudit(tenantID string, configDB configdb.Driver) (*AuditDetail, error) {
	ad := AuditDetail{}
	auditconf, err := configDB.GetAudit(tenantID)

	if err != nil {
		util.LogError("Error %v", err)
		return nil, err
	}

	ad.Enabled = auditconf.Enabled
	ad.TenantID = auditconf.TenantID

	return &ad, nil
}

//PutAudit changes the config for auditing for a given tenant.
//Inserts config database entry if one doesn't exist with defaults.
func PutAudit(tenantID string, configDB configdb.Driver) (*AuditDetail, error) {
	ad := AuditDetail{}
	auditconf, err := configDB.PutAudit(tenantID)

	if err != nil {
		util.LogError("Error %v", err)
		return nil, err
	}

	ad.Enabled = auditconf.Enabled

	return &ad, nil
}
