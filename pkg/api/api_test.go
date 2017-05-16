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
	"bytes"
	"net/http"
	"testing"

	"encoding/json"
	"github.com/databus23/goslo.policy"
	"github.com/sapcc/hermes/pkg/hermes"
	"github.com/sapcc/hermes/pkg/keystone"
	"github.com/sapcc/hermes/pkg/test"
	"github.com/spf13/viper"
	"io/ioutil"
)

type object map[string]interface{}

func setupTest(t *testing.T) http.Handler {
	// Initialise config for testing
	hermes.SetDefaultConfig()
	viper.SetConfigType("toml")
	var testConfigFile = []byte(`
[hermes]
storage_driver = "mock"
keystone_driver = "mock"
`)
	viper.ReadConfig(bytes.NewBuffer(testConfigFile))

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
	keystone := keystone.ConfiguredDriver()
	router, _ := NewV1Router(keystone)
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
		Path:             "/v1/events/d5eed458-6666-58ec-ad06-8d3cf6bafca1",
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
