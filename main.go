/*******************************************************************************
*
* Copyright 2022 SAP SE
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

package main

import (
	policy "github.com/databus23/goslo.policy"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
	"github.com/spf13/viper"

	"github.com/sapcc/hermes/internal/api"
	"github.com/sapcc/hermes/internal/identity"
	"github.com/sapcc/hermes/internal/storage"
	"github.com/sapcc/hermes/internal/util"
)

func main() {
	setDefaultConfig()
	bindEnvVariables()

	// Validate required Keystone authentication details
    if viper.GetString("keystone.username") == "" {
        logg.Fatal("Keystone username is not set")
    }
    if viper.GetString("keystone.password") == "" {
        logg.Fatal("Keystone password is not set")
    }
	if viper.GetString("elasticsearch.url") == "" {
		logg.Fatal("Elasticsearch URL is not set")
	}
	if viper.GetString("keystone.auth_url") == "" {
		logg.Fatal("Keystone authentication URL is not set")
	}

	logg.ShowDebug = viper.GetBool("hermes.debug")

	keystoneDriver := configuredKeystoneDriver()
	storageDriver := configuredStorageDriver()
	readPolicy()
	must.Succeed(api.Server(keystoneDriver, storageDriver))
}

func setDefaultConfig() {
	var nullEnforcer, err = policy.NewEnforcer(make(map[string]string))
	if err != nil {
		panic(err)
	}

	viper.SetDefault("hermes.debug", false)

	viper.SetDefault("hermes.keystone_driver", "keystone")
	viper.SetDefault("hermes.storage_driver", "elasticsearch")
	viper.SetDefault("hermes.PolicyEnforcer", &nullEnforcer)
	viper.SetDefault("hermes.PolicyFilePath", "/etc/policy.json")

	viper.SetDefault("elasticsearch.url", "")

	// Replace with your Keystone authentication URL
	viper.SetDefault("keystone.auth_url", "")

	// Replace with your Keystone authentication details
	viper.SetDefault("keystone.username", "")
	viper.SetDefault("keystone.password", "")
	viper.SetDefault("keystone.user_domain_name", "Default")
	viper.SetDefault("keystone.project_name", "service")
	viper.SetDefault("keystone.project_domain_name", "Default")
	viper.SetDefault("keystone.token_cache_time", 900)
	viper.SetDefault("keystone.memcached_servers", "")

	viper.SetDefault("API.ListenAddress", "0.0.0.0:8788")

	// index.max_result_window defaults to 10000, as per
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html
	// Increasing max_result_window to 20000, with corresponding changes to Elasticsearch to handle the increase.
	viper.SetDefault("elasticsearch.max_result_window", "20000")
}

// bindEnvVariables binds environment variables to viper keys
func bindEnvVariables() {
	must.Succeed(viper.BindEnv("hermes.keystone_driver", "HERMES_KEYSTONE_DRIVER"))
	must.Succeed(viper.BindEnv("hermes.storage_driver", "HERMES_STORAGE_DRIVER"))
	must.Succeed(viper.BindEnv("hermes.PolicyFilePath", "HERMES_POLICY_FILE_PATH"))

	must.Succeed(viper.BindEnv("elasticsearch.url", "HERMES_ES_URL"))

	must.Succeed(viper.BindEnv("keystone.auth_url", "HERMES_OS_AUTH_URL"))
	must.Succeed(viper.BindEnv("keystone.username", "HERMES_OS_USERNAME"))
	must.Succeed(viper.BindEnv("keystone.password", "HERMES_OS_PASSWORD"))
	must.Succeed(viper.BindEnv("keystone.user_domain_name", "HERMES_OS_USER_DOMAIN_NAME"))
	must.Succeed(viper.BindEnv("keystone.project_name", "HERMES_OS_PROJECT_NAME"))
	must.Succeed(viper.BindEnv("keystone.project_domain_name", "HERMES_OS_PROJECT_DOMAIN_NAME"))
	must.Succeed(viper.BindEnv("keystone.token_cache_time", "HERMES_OS_TOKEN_CACHE_TIME"))
	must.Succeed(viper.BindEnv("keystone.memcached_servers", "HERMES_OS_MEMCACHED_SERVERS"))

	must.Succeed(viper.BindEnv("API.ListenAddress", "HERMES_API_LISTEN_ADDRESS"))
	must.Succeed(viper.BindEnv("elasticsearch.username", "HERMES_ES_USERNAME"))
	must.Succeed(viper.BindEnv("elasticsearch.password", "HERMES_ES_PASSWORD"))
	must.Succeed(viper.BindEnv("elasticsearch.max_result_window", "HERMES_ES_MAX_RESULT_WINDOW"))
}

var keystoneIdentity = identity.Keystone{}
var mockIdentity = identity.Mock{}

func configuredKeystoneDriver() identity.Identity {
	driverName := viper.GetString("hermes.keystone_driver")
	switch driverName {
	case "keystone":
		return keystoneIdentity
	case "mock":
		return mockIdentity
	default:
		logg.Error("Couldn't match a keystone driver for configured value \"%s\"", driverName)
		return nil
	}
}

var elasticSearchStorage = storage.ElasticSearch{}
var mockStorage = storage.Mock{}

func configuredStorageDriver() storage.Storage {
	driverName := viper.GetString("hermes.storage_driver")
	switch driverName {
	case "elasticsearch":
		return elasticSearchStorage
	case "mock":
		return mockStorage
	default:
		logg.Error("Couldn't match a storage driver for configured value \"%s\"", driverName)
		return nil
	}
}

func readPolicy() {
	policyFilePath := viper.GetString("hermes.PolicyFilePath")
	if policyFilePath != "" {
		policyEnforcer, err := util.LoadPolicyFile(policyFilePath)
		if err != nil {
			logg.Fatal(err.Error())
		}
		if policyEnforcer != nil {
			viper.Set("hermes.PolicyEnforcer", policyEnforcer)
		}
	}
}
