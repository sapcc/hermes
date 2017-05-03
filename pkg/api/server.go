package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/hermes/pkg/keystone"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/sapcc/hermes/pkg/util"
	"github.com/spf13/viper"
)

func Server(keystone keystone.Interface, storage storage.Interface) error {
	fmt.Println("API")
	mainRouter := mux.NewRouter()

	//hook up the v1 API (this code is structured so that a newer API version can
	//be added easily later)
	v1Router, v1VersionData := NewV1Router(keystone)
	mainRouter.PathPrefix("/v1/").Handler(v1Router)

	//add the version advertisement that lists all available API versions
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allVersions := struct {
			Versions []VersionData `json:"versions"`
		}{[]VersionData{v1VersionData}}
		ReturnJSON(w, 300, allVersions)
	})

	//add Prometheus instrumentation
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", prometheus.InstrumentHandler("hermes-serve", mainRouter))

	//start HTTP server
	util.LogInfo("listening on " + viper.GetString("API.ListenAddress"))
	return http.ListenAndServe(viper.GetString("API.ListenAddress"), nil)
}
