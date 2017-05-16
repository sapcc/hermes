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
	"net/http"
	"os"

	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/databus23/goslo.policy"
	"github.com/sapcc/hermes/pkg/api"
	"github.com/sapcc/hermes/pkg/cmd"
	"github.com/sapcc/hermes/pkg/data"
	"github.com/sapcc/hermes/pkg/hermes"
	"github.com/sapcc/hermes/pkg/keystone"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
)

var configPath *string

func main() {
	if os.Getenv("HERMES_INSECURE") == "1" {
		fmt.Println("Insecure HTTPS mode!")
		http.DefaultClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	parseCmdlineFlags()

	hermes.SetDefaultConfig()
	readConfig(configPath)
	readPolicy()

	// If there are args left over after flag processing, we are a Hermes CLI client
	if len(flag.Args()) > 0 {
		cmd.RootCmd.SetArgs(flag.Args())
		if err := cmd.RootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else { // otherwise, we are running a Hermes API server
		api.Server(keystone.ConfiguredDriver(), storage.ConfiguredDriver())
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
	for _, os_var_name := range data.OS_vars {
		viper.BindEnv("keystone."+os_var_name, "OS_"+strings.ToUpper(os_var_name))
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
