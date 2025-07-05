// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/respondwith"
)

// ReturnESJSON is a custom response helper that preserves Elasticsearch URL formatting.
// This is needed because Elasticsearch requires literal & characters in URLs, but Go's
// JSON marshaler escapes them as \u0026.
//
// Example:
//
//	events := []Event{...}
//	ReturnESJSON(w, http.StatusOK, map[string]any{
//		"events": events,
//		"total":  len(events),
//	})
func ReturnESJSON(w http.ResponseWriter, code int, data any) {
	payload, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		respondwith.ErrorText(w, err)
		return
	}

	// Replace escaped ampersands with literal ones for Elasticsearch compatibility
	payload = bytes.ReplaceAll(payload, []byte("\\u0026"), []byte("&"))

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
//
// Example:
//
//	events, err := storage.GetEvents(filter)
//	if ReturnError(w, err) {
//		return
//	}
func ReturnError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
	return true
}

// getProtocol determines the protocol (http or https) for building URLs.
func getProtocol(req *http.Request) string {
	protocol := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	return protocol
}

// ValidationError writes a 400 Bad Request error response following go-bits patterns.
func ValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, writeErr := w.Write([]byte(err.Error()))
	if writeErr != nil {
		logg.Error("Failed to write validation error response: %s", writeErr.Error())
	}
}

// NotFoundError writes a 404 Not Found error response following go-bits patterns.
func NotFoundError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, writeErr := w.Write([]byte(err.Error()))
	if writeErr != nil {
		logg.Error("Failed to write not found error response: %s", writeErr.Error())
	}
}
