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

package configdb

// Driver is an interface that wraps the underlying event storage mechanism.
// Because it is an interface, the real implementation can be mocked away in unit tests.
type Driver interface {
	/********** requests to MySQL **********/
	GetAudit(tenantId string) (*AuditConfig, error)
	PutAudit(tenantId string) (*AuditConfig, error)
}

// AuditConfig contains the mapping to MySQL config table.
type AuditConfig struct {
	Enabled  bool   `db:"enabled"`
	TenantID string `db:"tenant_id, primarykey"`
}
