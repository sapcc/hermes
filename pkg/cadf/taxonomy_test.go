package cadf

import "testing"

func TestIsTypeURI(t *testing.T) {
	type args struct {
		TypeURI string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{TypeURI: "storage"}, true},
		{"test2", args{TypeURI: "storage/data"}, true},
		{"test3", args{TypeURI: "unknown"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTypeURI(tt.args.TypeURI); got != tt.want {
				t.Errorf("IsTypeURI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAction(t *testing.T) {
	type args struct {
		Action string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{Action: "create"}, true},
		{"test2", args{Action: "delete"}, true},
		{"test3", args{Action: "bork"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAction(tt.args.Action); got != tt.want {
				t.Errorf("IsAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsOutcome(t *testing.T) {
	type args struct {
		outcome string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{outcome: "success"}, true},
		{"test2", args{outcome: "failure"}, true},
		{"test3", args{outcome: "bork"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOutcome(tt.args.outcome); got != tt.want {
				t.Errorf("IsOutcome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAction(t *testing.T) {
	type args struct {
		req string
	}
	tests := []struct {
		name       string
		args       args
		wantAction string
	}{
		{"test1", args{req: "get"}, "read"},
		{"test2", args{req: "post"}, "create"},
		{"test3", args{req: "bork"}, "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotAction := GetAction(tt.args.req); gotAction != tt.wantAction {
				t.Errorf("GetAction() = %v, want %v", gotAction, tt.wantAction)
			}
		})
	}
}
