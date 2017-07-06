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

package api

import (
	"fmt"
	"github.com/sapcc/hermes/pkg/hermes"
	"net/http"
)

//GetAudit handles GET /v1/audit.
//Check permissions for actions
func (p *v1Provider) GetAudit(res http.ResponseWriter, req *http.Request) {
	// QueryString - TenantId
	token := p.CheckToken(req)
	//if !token.Require(res, "audit:show") {
	//	//return
	//	//TODO - This is supposed to work, and isn't. So FIXME
	//}

	tenantId, err := getTenantId(token, req, res)
	if ReturnError(res, err) {
		return
	}

	auditconf, err := hermes.GetAudit(tenantId, p.configdb)

	if ReturnError(res, err) {
		return
	}
	if auditconf == nil {
		err := fmt.Errorf("Audit Configuration could not be found for tenant %s", tenantId)
		http.Error(res, err.Error(), 404)
		return
	}
	ReturnJSON(res, 200, auditconf)
}

//PutAudit handles PUT /v1/audit.
func (p *v1Provider) PutAudit(res http.ResponseWriter, req *http.Request) {
	// Check Authorizations
	token := p.CheckToken(req)
	//if !token.Require(res, "audit:edit") {
	//	//return
	////	//TODO - This is supposed to work, and isn't. So FIXME
	//}

	tenantId, err := getTenantId(token, req, res)
	if ReturnError(res, err) {
		return
	}

	auditconf, err := hermes.PutAudit(tenantId, p.configdb)

	if auditconf == nil {
		err := fmt.Errorf("Audit Configuration could not be found for tenant %s", tenantId)
		http.Error(res, err.Error(), 404)
		return
	}
	ReturnJSON(res, 200, auditconf)
}

//Check all the input, make sure they have permission to do it,
// call the business logic, then convert to  json.
