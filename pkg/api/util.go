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

package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/hermes/pkg/util"
)

// utility functionality

//VersionData is used by version advertisement handlers.
type VersionData struct {
	Status string            `json:"status"`
	ID     string            `json:"id"`
	Links  []versionLinkData `json:"links"`
}

//versionLinkData is used by version advertisement handlers, as part of the
//VersionData struct.
type versionLinkData struct {
	URL      string `json:"href"`
	Relation string `json:"rel"`
	Type     string `json:"type,omitempty"`
}

//ReturnJSON is a convenience function for HTTP handlers returning JSON data.
//The `code` argument specifies the HTTP response code, usually 200.
func ReturnJSON(w http.ResponseWriter, code int, data interface{}) {
	payload, err := json.MarshalIndent(&data, "", "  ")
	// Replaces & symbols properly in json within urls due to Elasticsearch
	payload = bytes.Replace(payload, []byte("\\u0026"), []byte("&"), -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(payload)
	if err != nil {
		util.LogDebug("Issue with writing payload when returning Json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//ReturnError produces an error response with HTTP status code 500 if the given
//error is non-nil. Otherwise, nothing is done and false is returned.
func ReturnError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	http.Error(w, err.Error(), 500)
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
	durationSummary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{Name: "hermes_request_duration_seconds", Help: "Duration/latency of a Hermes request", ConstLabels: prometheus.Labels{"handler": handler}}, nil)
	prometheus.MustRegister(durationSummary)

	return promhttp.InstrumentHandlerDuration(durationSummary, handlerFunc)
}

func observeResponseSize(handlerFunc http.HandlerFunc, handler string) http.HandlerFunc {
	durationSummary := prometheus.NewSummaryVec(prometheus.SummaryOpts{Name: "hermes_response_size_bytes", Help: "Size of the Hermes response (e.g. to a query)", ConstLabels: prometheus.Labels{"handler": handler}}, nil)
	prometheus.MustRegister(durationSummary)

	return promhttp.InstrumentHandlerResponseSize(durationSummary, http.HandlerFunc(handlerFunc)).ServeHTTP
}
