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
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/databus23/goslo.policy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/hermes/pkg/identity"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/sapcc/hermes/pkg/test"
	"github.com/spf13/viper"
)

type object map[string]interface{}

func setupTest(t *testing.T) http.Handler {
	//load test policy (where everything is allowed)
	policyBytes, err := ioutil.ReadFile("../test/policy.json")
	if err != nil {
		t.Fatal(err)
	}
	policyRules := make(map[string]string)
	err = json.Unmarshal(policyBytes, &policyRules)
	if err != nil {
		t.Fatal(err)
	}
	policyEnforcer, err := policy.NewEnforcer(policyRules)
	if err != nil {
		t.Fatal(err)
	}
	viper.Set("hermes.PolicyEnforcer", policyEnforcer)

	//create test driver with the domains and projects from start-data.sql
	keystone := identity.Mock{}
	storage := storage.Mock{}

	prometheus.DefaultRegisterer = prometheus.NewPedanticRegistry()
	router, _ := NewV1Handler(keystone, storage)
	return router
}

func Test_API(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		path       string
		statuscode int
		json       string
	}{
		{"Metadata", "GET", "/v1/", http.StatusOK, "fixtures/api-metadata.json"},
		{"EventDetails", "GET", "/v1/events/7be6c4ff-b761-5f1f-b234-f5d41616c2cd", http.StatusOK, "fixtures/event-details.json"},
		{"EventList", "GET", "/v1/events?event_type=identity.project.deleted&offset=10", http.StatusOK, "fixtures/event-list.json"},
		{"Attributes", "GET", "/v1/attributes/resource_type", http.StatusOK, "fixtures/attributes.json"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := setupTest(t)

			test.APIRequest{
				Method:           tc.method,
				Path:             tc.path,
				ExpectStatusCode: tc.statuscode,
				ExpectJSON:       tc.json,
			}.Check(t, router)
		})
	}

}
