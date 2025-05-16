// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/spf13/viper"

	"github.com/sapcc/go-bits/gopherpolicy"
	"github.com/sapcc/go-bits/httpext"
	"github.com/sapcc/go-bits/logg"

	"github.com/sapcc/hermes/pkg/storage"
)

// Server Set up and start the API server, hooking it up to the API router
func Server(validator gopherpolicy.Validator, storageInterface storage.Storage) error {
	fmt.Println("API")
	mainRouter := setupRouter(validator, storageInterface)

	// start HTTP server
	listenaddress := viper.GetString("API.ListenAddress")
	logg.Info("listening on %s", listenaddress)
	// enable cors support
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"X-Auth-Token", "Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "HEAD"},
		MaxAge:         600,
	})
	handler := c.Handler(mainRouter)

	ctx := httpext.ContextWithSIGINT(context.Background(), 10*time.Second)
	return httpext.ListenAndServeContext(ctx, listenaddress, handler)
}

func setupRouter(validator gopherpolicy.Validator, storageInterface storage.Storage) http.Handler {
	mainRouter := mux.NewRouter()
	// hook up the v1 API (this code is structured so that a newer API version can
	// be added easily later)
	v1Router, v1VersionData := NewV1Handler(validator, storageInterface)
	mainRouter.PathPrefix("/v1/").Handler(v1Router)

	// add the version advertisement that lists all available API versions
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		allVersions := struct {
			Versions []VersionData `json:"versions"`
		}{[]VersionData{v1VersionData}}
		ReturnJSON(w, http.StatusMultipleChoices, allVersions)
	})

	// instrumentation
	mainRouter.Handle("/metrics", promhttp.Handler())

	return gaugeInflight(mainRouter)
}
