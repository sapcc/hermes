// Copyright 2022 SAP SE
// SPDX-FileCopyrightText: 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"crypto/tls"
	"net/http"
	"os"
)

func init() {
	// I have some trouble getting hermes to connect to our staging OpenStack
	// through mitmproxy (which is very useful for development and debugging) when
	// TLS certificate verification is enabled. Therefore, allow to turn it off
	// with an env variable. (It's very important that this is not the standard
	// "DEBUG" variable. "DEBUG" is meant to be useful for production systems,
	// where you definitely don't want to turn off certificate verification.)
	if os.Getenv("HERMES_INSECURE") == "1" {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // intentional usage of InsecureSkipVerify
		}
		http.DefaultTransport = &http.Transport{
			TLSClientConfig: tlsConf,
			Proxy:           http.ProxyFromEnvironment,
		}
	}
}
