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

	"github.com/sapcc/hermes/pkg/api"
	"github.com/sapcc/hermes/pkg/cli"
	"github.com/spf13/viper"
)

func main() {
	// Get config file location
	configPath := flag.String("f", "hermes.conf", "specifies the location of the TOML-format configuration file")
	flag.Usage = printUsage
	flag.Parse()

	//Don't read config file if the default config file isn't there, as we will just use
	// config defaults in that case
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

	// If there are args left over after flag processing, we are a Hermes CLI client
	if len(flag.Args()) > 0 {
		cli.Command(flag.Args())
	} else { // otherwise, we are running a Hermes API server
		api.Server()
	}
	//fmt.Println("Selected driver:", viper.Get("hermes.driver"))
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

////////////////////////////////////////////////////////////////////////////////
// task: serve
/*
func taskServe(config hermes.Configuration, driver hermes.Driver, args []string) error {
	if len(args) != 0 {
		printUsageAndExit()
	}

	mainRouter := mux.NewRouter()

	//hook up the v1 API (this code is structured so that a newer API version can
	//be added easily later)
	v1Router, v1VersionData := api.NewV1Router(driver, config)
	mainRouter.PathPrefix("/v1/").Handler(v1Router)

	//add the version advertisement that lists all available API versions
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allVersions := struct {
			Versions []api.VersionData `json:"versions"`
		}{[]api.VersionData{v1VersionData}}
		api.ReturnJSON(w, 300, allVersions)
	})

	//add Prometheus instrumentation
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", prometheus.InstrumentHandler("hermes-serve", mainRouter))

	//start HTTP server
	util.LogInfo("listening on " + config.API.ListenAddress)
	return http.ListenAndServe(config.API.ListenAddress, nil)
}
*/
