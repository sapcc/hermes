package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
