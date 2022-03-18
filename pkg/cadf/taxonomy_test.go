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
