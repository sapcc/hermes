package configdb

//Mock implementation of configdb with static data
type Mock struct{}

//GetAudit Mock implementation of GetAudit Endpoint
func (m Mock) GetAudit(tenantID string) (*AuditConfig, error) {
	d := AuditConfig{Enabled: true}

	return &d, nil
}

//PutAudit Mock implementation of PutAudit Endpoint
func (m Mock) PutAudit(tenantID string) (*AuditConfig, error) {
	d := AuditConfig{Enabled: true}

	return &d, nil
}
