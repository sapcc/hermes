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
	"github.com/gorilla/mux"
	"github.com/sapcc/hermes/pkg/identity"
	"github.com/sapcc/hermes/pkg/storage"
	"net/http"
	"strings"
)

type v1Provider struct {
	keystone    identity.Identity
	storage     storage.Storage
	versionData VersionData
}

//NewV1Handler creates a http.Handler that serves the Hermes v1 API.
//It also returns the VersionData for this API version which is needed for the
//version advertisement on "GET /".
func NewV1Handler(keystone identity.Identity, storage storage.Storage) (http.Handler, VersionData) {
	r := mux.NewRouter()

	p := &v1Provider{
		keystone: keystone,
		storage:  storage,
	}
	p.versionData = VersionData{
		Status: "CURRENT",
		ID:     "v1",
		Links: []versionLinkData{
			{
				Relation: "self",
				URL:      p.Path(),
			},
			{
				Relation: "describedby",
				URL:      "https://github.com/sapcc/hermes/tree/master/docs",
				Type:     "text/html",
			},
		},
	}

	r.Methods("GET").Path("/v1/").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		p.versionData.Links[0].URL = fmt.Sprintf("%s://%s%s/", getProtocol(req), req.Host, p.Path())
		ReturnJSON(res, 200, map[string]interface{}{"version": p.versionData})
	})

	r.Methods("GET").Path("/v1/events").HandlerFunc(
		observeDuration(observeResponseSize(p.ListEvents, "ListEvents"), "ListEvents"))
	r.Methods("GET").Path("/v1/events/{event_id}").HandlerFunc(
		observeDuration(observeResponseSize(p.GetEventDetails, "GetEventDetails"), "GetEventDetails"))
	r.Methods("GET").Path("/v1/attributes/{attribute_name}").HandlerFunc(
		observeDuration(observeResponseSize(p.GetAttributes, "GetAttributes"), "GetAttributes"))

	return r, p.versionData
}

//Path constructs a full URL for a given URL path below the /v1/ endpoint.
func (p *v1Provider) Path(elements ...string) string {
	parts := []string{
		strings.TrimSuffix( /*p.Driver.Cluster().Config.CatalogURL*/ "", "/"),
		"v1",
	}
	parts = append(parts, elements...)
	return strings.Join(parts, "/")
}
