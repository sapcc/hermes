// Copyright 2022 SAP SE
// SPDX-FileCopyrightText: 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/go-bits/logg"
)

// utility functionality

// VersionData is used by version advertisement handlers.
type VersionData struct {
	Status string            `json:"status"`
	ID     string            `json:"id"`
	Links  []versionLinkData `json:"links"`
}

// versionLinkData is used by version advertisement handlers, as part of the
// VersionData struct.
type versionLinkData struct {
	URL      string `json:"href"`
	Relation string `json:"rel"`
	Type     string `json:"type,omitempty"`
}

// ReturnJSON is a convenience function for HTTP handlers returning JSON data.
// The `code` argument specifies the HTTP response code, usually 200.
func ReturnJSON(w http.ResponseWriter, code int, data any) {
	payload, err := json.MarshalIndent(&data, "", "  ")
	// Replaces & symbols properly in json within urls due to Elasticsearch
	payload = bytes.ReplaceAll(payload, []byte("\\u0026"), []byte("&"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(payload)
	if err != nil {
		// It's too late to write this as a 5xx response since we've already
		// started writing a 2xx response, so this can only be logged.
		logg.Error("Issue with writing payload when returning JSON: %s", err.Error())
	}
}

// ReturnError produces an error response with HTTP status code 500 if the given
// error is non-nil. Otherwise, nothing is done and false is returned.
func ReturnError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
	return true
}

var authErrorsCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "hermes_logon_errors_count", Help: "Number of logon errors occurred"})
var authFailuresCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "hermes_logon_failures_count", Help: "Number of logon attempts failed due to wrong credentials"})
var storageErrorsCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "hermes_storage_errors_count", Help: "Number of technical errors occurred when accessing underlying storage (i.e. Elasticsearch)"})

func init() {
	prometheus.MustRegister(authErrorsCounter, authFailuresCounter, storageErrorsCounter)
}

func gaugeInflight(handler http.Handler) http.Handler {
	inflightGauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: "hermes_requests_inflight", Help: "Number of inflight HTTP requests served by Hermes"})
	prometheus.MustRegister(inflightGauge)

	return promhttp.InstrumentHandlerInFlight(inflightGauge, handler)
}

func observeDuration(handlerFunc http.HandlerFunc, handler string) http.HandlerFunc {
	durationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "hermes_request_duration_seconds",
			Help:        "Duration/latency of a Hermes request",
			ConstLabels: prometheus.Labels{"handler": handler},
			Buckets:     prometheus.DefBuckets,
		},
		[]string{},
	)
	prometheus.MustRegister(durationHistogram)

	return promhttp.InstrumentHandlerDuration(durationHistogram, handlerFunc)
}

func observeResponseSize(handlerFunc http.HandlerFunc, handler string) http.HandlerFunc {
	responseSizeHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "hermes_response_size_bytes",
			Help:        "Size of the Hermes response (e.g. to a query)",
			ConstLabels: prometheus.Labels{"handler": handler},
			Buckets:     prometheus.LinearBuckets(100, 100, 10),
		},
		[]string{},
	)
	prometheus.MustRegister(responseSizeHistogram)

	return promhttp.InstrumentHandlerResponseSize(responseSizeHistogram, handlerFunc).ServeHTTP
}
