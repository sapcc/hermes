package util

import (
	"encoding/json"
	"io/ioutil"

	policy "github.com/databus23/goslo.policy"
)

// LoadPolicyFile used to Load the hermes policy.json file from disk.
func LoadPolicyFile(path string) (*policy.Enforcer, error) {
	if path == "" {
		return nil, nil
	}
	bytes, err := ioutil.ReadFile(path)
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
