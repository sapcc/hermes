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
	"net/http"
	"testing"

	"encoding/json"
	"github.com/databus23/goslo.policy"
	"github.com/sapcc/hermes/pkg/configdb"
	"github.com/sapcc/hermes/pkg/identity"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/sapcc/hermes/pkg/test"
	"github.com/spf13/viper"
	"io/ioutil"
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
	configdb := configdb.Mock{}
	router, _ := NewV1Router(keystone, storage, configdb)
	return router
}

func Test_APIMetadata(t *testing.T) {
	router := setupTest(t)

	test.APIRequest{
		Method:           "GET",
		Path:             "/v1/",
		ExpectStatusCode: 200,
		ExpectJSON:       "fixtures/api-metadata.json",
	}.Check(t, router)

}

func Test_APIGetEventDetails(t *testing.T) {
	router := setupTest(t)

	test.APIRequest{
		Method:           "GET",
		Path:             "/v1/events/7be6c4ff-b761-5f1f-b234-f5d41616c2cd",
		ExpectStatusCode: 200,
		ExpectJSON:       "fixtures/event-details.json",
	}.Check(t, router)

}

func Test_APIGetEventList(t *testing.T) {
	router := setupTest(t)

	test.APIRequest{
		Method:           "GET",
		Path:             "/v1/events?event_type=identity.project.deleted&offset=10",
		ExpectStatusCode: 200,
		ExpectJSON:       "fixtures/event-list.json",
	}.Check(t, router)

}
