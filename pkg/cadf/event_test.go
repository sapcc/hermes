package cadf

import (
	"testing"
)

func TestStripPort(t *testing.T) {
	type args struct {
		hostPort string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"testLocalhost", args{"localhost:8080"}, "localhost"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripPort(tt.args.hostPort); got != tt.want {
				t.Errorf("StripPort() = %v, want %v", got, tt.want)
			}
		})
	}
}
