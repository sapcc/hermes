// Copyright 2017 SAP SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"fmt"

	"encoding/json"
	"github.com/sapcc/hermes/pkg/hermes"
	"github.com/spf13/cobra"
	"github.com/sapcc/hermes/pkg/cmd/auth"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get the audit event with ID <id>",
	Long:  `Get the audit event with ID <id>.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("You must specify exactly one event ID.")
		}

		token := auth.GetToken(keystoneDriver)
		if !token.Require("event:show") {
			return errors.New("You are not authorised to view event details")
		}

		eventId := args[0]
		event, err := hermes.GetEvent(eventId, token.TenantId(), keystoneDriver, storageDriver)
		if err != nil {
			return err
		}
		if event == nil {
			return errors.New(fmt.Sprintf("Event %s could not be found in tenant %s", eventId, token.TenantId()))
		}
		json, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			return err
		}
		fmt.Printf("%s", json)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
