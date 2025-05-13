// Copyright 2022 SAP SE
// SPDX-FileCopyrightText: 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"encoding/json"
	"os"

	policy "github.com/databus23/goslo.policy"
)

// LoadPolicyFile used to Load the hermes policy.json file from disk.
func LoadPolicyFile(path string) (*policy.Enforcer, error) {
	if path == "" {
		return nil, nil
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules map[string]string
	err = json.Unmarshal(bytes, &rules)
	if err != nil {
		return nil, err
	}
	return policy.NewEnforcer(rules)
}
