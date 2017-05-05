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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sapcc/hermes/pkg/keystone"
)

//VersionData is used by version advertisement handlers.
type VersionData struct {
	Status string            `json:"status"`
	ID     string            `json:"id"`
	Links  []VersionLinkData `json:"links"`
}

//VersionLinkData is used by version advertisement handlers, as part of the
//VersionData struct.
type VersionLinkData struct {
	URL      string `json:"href"`
	Relation string `json:"rel"`
	Type     string `json:"type,omitempty"`
}

type v1Provider struct {
	Keystone    keystone.Interface
	VersionData VersionData
}

//NewV1Router creates a http.Handler that serves the Limes v1 API.
//It also returns the VersionData for this API version which is needed for the
//version advertisement on "GET /".
func NewV1Router(keystone keystone.Interface) (http.Handler, VersionData) {
	r := mux.NewRouter()
	p := &v1Provider{
		Keystone: keystone,
	}
	p.VersionData = VersionData{
		Status: "CURRENT",
		ID:     "v1",
		Links: []VersionLinkData{
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

	r.Methods("GET").Path("/v1/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ReturnJSON(w, 200, map[string]interface{}{"version": p.VersionData})
	})

	//r.Methods("GET").Path("/v1/clusters").HandlerFunc(p.ListClusters)
	//r.Methods("GET").Path("/v1/clusters/{cluster_id}").HandlerFunc(p.GetCluster)
	//r.Methods("PUT").Path("/v1/clusters/{cluster_id}").HandlerFunc(p.PutCluster)
	//
	//r.Methods("GET").Path("/v1/domains").HandlerFunc(p.ListDomains)
	//r.Methods("GET").Path("/v1/domains/{domain_id}").HandlerFunc(p.GetDomain)
	//r.Methods("POST").Path("/v1/domains/discover").HandlerFunc(p.DiscoverDomains)
	//r.Methods("PUT").Path("/v1/domains/{domain_id}").HandlerFunc(p.PutDomain)
	//
	//r.Methods("GET").Path("/v1/domains/{domain_id}/projects").HandlerFunc(p.ListProjects)
	//r.Methods("GET").Path("/v1/domains/{domain_id}/projects/{project_id}").HandlerFunc(p.GetProject)
	//r.Methods("POST").Path("/v1/domains/{domain_id}/projects/discover").HandlerFunc(p.DiscoverProjects)
	//r.Methods("POST").Path("/v1/domains/{domain_id}/projects/{project_id}/sync").HandlerFunc(p.SyncProject)
	//r.Methods("PUT").Path("/v1/domains/{domain_id}/projects/{project_id}").HandlerFunc(p.PutProject)

	return r, p.VersionData
}

//ReturnJSON is a convenience function for HTTP handlers returning JSON data.
//The `code` argument specifies the HTTP response code, usually 200.
func ReturnJSON(w http.ResponseWriter, code int, data interface{}) {
	bytes, err := json.MarshalIndent(&data, "", "  ")
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(bytes)
	} else {
		http.Error(w, err.Error(), 500)
	}
}

//ReturnError produces an error response with HTTP status code 500 if the given
//error is non-nil. Otherwise, nothing is done and false is returned.
func ReturnError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	http.Error(w, err.Error(), 500)
	return true
}

//RequireJSON will parse the request body into the given data structure, or
//write an error response if that fails.
func RequireJSON(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		http.Error(w, "request body is not valid JSON: "+err.Error(), 400)
		return false
	}
	return true
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
