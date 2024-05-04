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

package storage

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "No duplicates",
			input:    []string{"apple", "banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "With duplicates",
			input:    []string{"apple", "banana", "cherry", "apple", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "All duplicates",
			input:    []string{"apple", "apple", "apple"},
			expected: []string{"apple"},
		},
		{
			name:     "Single element",
			input:    []string{"apple"},
			expected: []string{"apple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := RemoveDuplicates(tt.input)
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("got %v, want %v", output, tt.expected)
			}
		})
	}
}
