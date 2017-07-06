package configdb

type Mock struct{}

// Mock configdb driver with static data

func (m Mock) GetAudit(tenantId string) (*AuditConfig, error) {
	d := AuditConfig{Enabled: true}

	return &d, nil
}

func (m Mock) PutAudit(tenantId string) (*AuditConfig, error) {
	d := AuditConfig{Enabled: true}

	return &d, nil
}
