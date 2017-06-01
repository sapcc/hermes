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

package main

import (
	"flag"
	"fmt"
	"os"

	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/databus23/goslo.policy"
	"github.com/sapcc/hermes/pkg/api"
	"github.com/sapcc/hermes/pkg/cmd"
	"github.com/sapcc/hermes/pkg/keystone"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
	"log"
)

var configPath *string

func main() {
	parseCmdlineFlags()

	setDefaultConfig()
	readConfig(configPath)
	keystoneDriver := configuredKeystoneDriver()
	storageDriver := configuredStorageDriver()
	readPolicy()

	// If there are args left over after flag processing, we are a Hermes CLI client
	if len(flag.Args()) > 0 {
		cmd.RootCmd.SetArgs(flag.Args())
		cmd.SetDrivers(keystoneDriver, storageDriver)
		if err := cmd.RootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else { // otherwise, we are running a Hermes API server
		api.Server(keystoneDriver, storageDriver)
	}
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
	viper.SetDefault("hermes.keystone_driver", "keystone")
	viper.SetDefault("hermes.storage_driver", "elasticsearch")
	viper.SetDefault("hermes.enrich_keystone_events", "False")
	viper.SetDefault("API.ListenAddress", "0.0.0.0:8788")
	viper.SetDefault("elasticsearch.url", "localhost:9200")
	// index.max_result_window defaults to 10000, as per
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules.html
	viper.SetDefault("elasticsearch.max_result_window", "10000")
}

func readConfig(configPath *string) {
	// Don't read config file if the default config file isn't there,
	//  as we will just fall back to config defaults in that case
	var shouldReadConfig = true
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		shouldReadConfig = *configPath != flag.Lookup("f").DefValue
	}
	// Now we sorted that out, read the config
	if shouldReadConfig {
		viper.SetConfigFile(*configPath)
		viper.SetConfigType("toml")
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	// Setup environment variable overrides for OpenStack authentication
	for _, osVarName := range cmd.OSVars {
		viper.BindEnv("keystone."+osVarName, "OS_"+strings.ToUpper(osVarName))
	}

}

func configuredKeystoneDriver() keystone.Driver {
	driverName := viper.GetString("hermes.keystone_driver")
	switch driverName {
	case "keystone":
		return keystone.Keystone()
	case "mock":
		return keystone.Mock()
	default:
		log.Printf("Couldn't match a keystone driver for configured value \"%s\"", driverName)
		return nil
	}
}

func configuredStorageDriver() storage.Driver {
	driverName := viper.GetString("hermes.storage_driver")
	switch driverName {
	case "elasticsearch":
		return storage.ElasticSearch()
	case "mock":
		return storage.Mock()
	default:
		log.Printf("Couldn't match a storage driver for configured value \"%s\"", driverName)
		return nil
	}
}

func readPolicy() {
	//load the policy file
	policyEnforcer, err := loadPolicyFile(viper.GetString("hermes.PolicyFilePath"))
	if err != nil {
		util.LogFatal(err.Error())
	}
	viper.Set("hermes.PolicyEnforcer", policyEnforcer)
}

func loadPolicyFile(path string) (*policy.Enforcer, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules map[string]string
	err = json.Unmarshal(bytes, &rules)
	if err != nil {
		return nil, err
	}
	return policy.NewEnforcer(rules)
}
