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
	"flag"
	"fmt"
	"os"

	policy "github.com/databus23/goslo.policy"
	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/must"
	"github.com/sapcc/go-bits/osext"
	"github.com/spf13/viper"

	"github.com/sapcc/hermes/internal/api"
	"github.com/sapcc/hermes/internal/identity"
	"github.com/sapcc/hermes/internal/storage"
	"github.com/sapcc/hermes/internal/util"
)

var configPath *string

func main() {
	logg.ShowDebug = osext.GetenvBool("HERMES_DEBUG")
	parseCmdlineFlags()

	setDefaultConfig()
	readConfig(configPath)
	keystoneDriver := configuredKeystoneDriver()
	storageDriver := configuredStorageDriver()
	readPolicy()
	must.Succeed(api.Server(keystoneDriver, storageDriver))
}

func parseCmdlineFlags() {
	// Get config file location
	configPath = flag.String("f", "hermes.conf", "specifies the location of the TOML-format configuration file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func setDefaultConfig() {
	var nullEnforcer, err = policy.NewEnforcer(make(map[string]string))
	if err != nil {
		panic(err)
	}
	viper.SetDefault("hermes.keystone_driver", "keystone")
	viper.SetDefault("hermes.storage_driver", "elasticsearch")
	viper.SetDefault("hermes.PolicyEnforcer", &nullEnforcer)
	viper.SetDefault("API.ListenAddress", "0.0.0.0:8788")
	viper.SetDefault("elasticsearch.url", "localhost:9200")
	// index.max_result_window defaults to 10000, as per
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html
	// Increasing max_result_window to 20000, with corresponding changes to Elasticsearch to handle the increase.
	viper.SetDefault("elasticsearch.max_result_window", "20000")
}

func readConfig(configPath *string) {
	// Don't read config file if the default config file isn't there,
	//  as we will just fall back to config defaults in that case
	var shouldReadConfig = true
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		shouldReadConfig = *configPath != flag.Lookup("f").DefValue
	}
	// Now we sorted that out, read the config
	logg.Debug("Should read config: %v, config file is %s", shouldReadConfig, *configPath)
	if shouldReadConfig {
		viper.SetConfigFile(*configPath)
		viper.SetConfigType("toml")
		must.Succeed(viper.ReadInConfig())
	}
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
	//load the policy file
	policyEnforcer, err := util.LoadPolicyFile(viper.GetString("hermes.PolicyFilePath"))
	if err != nil {
		logg.Fatal(err.Error())
	}
	if policyEnforcer != nil {
		viper.Set("hermes.PolicyEnforcer", policyEnforcer)
	}
}
