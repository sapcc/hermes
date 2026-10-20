// SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/sapcc/hermes/pkg/identity"
)

func TestConfigureDataplaneWriter(t *testing.T) {
	svc := identity.ServiceUser{ID: "svc-id", Name: "hermes", Domain: "Default"}
	tests := []struct {
		name       string
		configured string
		svc        identity.ServiceUser
		want       string
	}{
		{"configured value wins over the service user", "cfg-id", svc, "cfg-id"},
		{"service user is used when unset", "", svc, "svc-id"},
		{"whitespace-only counts as unset", "  ", svc, "svc-id"},
		{"configured value is trimmed", " cfg-id ", svc, "cfg-id"},
		{"whitespace-only and no service user ends empty", " ", identity.ServiceUser{}, ""},
		{"nothing available stays empty", "", identity.ServiceUser{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(viper.Reset)
			viper.Set("hermes.dataplane_writer_user_id", tt.configured)
			configureDataplaneWriter(tt.svc)
			if got := viper.GetString("hermes.dataplane_writer_user_id"); got != tt.want {
				t.Errorf("viper has %q, want %q", got, tt.want)
			}
		})
	}
}
