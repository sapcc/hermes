// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

func clientWithAuthResult(t *testing.T, body any) *gophercloud.ProviderClient {
	t.Helper()
	var result tokens.CreateResult
	result.Body = body
	result.Header = http.Header{"X-Subject-Token": []string{"tok"}}
	client := &gophercloud.ProviderClient{}
	if err := client.SetTokenAndAuthResult(result); err != nil {
		t.Fatal(err)
	}
	return client
}

func TestServiceUserID(t *testing.T) {
	withUser := map[string]any{"token": map[string]any{"user": map[string]any{"id": "user-123", "name": "hermes"}}}
	noUser := map[string]any{"token": map[string]any{}}
	emptyID := map[string]any{"token": map[string]any{"user": map[string]any{"name": "hermes"}}}

	tests := []struct {
		name    string
		client  *gophercloud.ProviderClient
		want    string
		wantErr bool
	}{
		{"user in token", clientWithAuthResult(t, withUser), "user-123", false},
		{"token without user", clientWithAuthResult(t, noUser), "", true},
		{"user without id", clientWithAuthResult(t, emptyID), "", true},
		{"no auth result", &gophercloud.ProviderClient{}, "", true},
		{"nil client", nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ServiceUserID(tt.client)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveWriterUserID(t *testing.T) {
	tests := []struct {
		name, configured, service, wantID, wantSource string
	}{
		{"configured wins", "cfg", "svc", "cfg", WriterSourceConfig},
		{"empty config uses service user", "", "svc", "svc", WriterSourceServiceToken},
		{"neither stays empty", "", "", "", WriterSourceServiceToken},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, source := ResolveWriterUserID(tt.configured, tt.service)
			if id != tt.wantID || source != tt.wantSource {
				t.Errorf("got (%q, %q), want (%q, %q)", id, source, tt.wantID, tt.wantSource)
			}
		})
	}
}
