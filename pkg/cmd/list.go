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
	"github.com/olekukonko/tablewriter"
	"github.com/sapcc/hermes/pkg/cmd/auth"
	"github.com/sapcc/hermes/pkg/data"
	"github.com/sapcc/hermes/pkg/hermes"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists a project’s or domain's audit events. The project or domain comes from the scope of the authentication parameters.",
	Long: `Lists a project’s or domain's audit events. The project or domain comes from the scope of the authentication parameters.

Date Filters:
The value for the time parameter is a comma-separated list of time stamps in ISO 8601 format. The time stamps can be prefixed with any of these comparison operators: gt: (greater-than), gte: (greater-than-or-equal), lt: (less-than), lte: (less-than-or-equal).
For example, to get a list of events that will expire in January of 2020:
GET /v1/events?time=gte:2020-01-01T00:00:00,lt:2020-02-01T00:00:00

Sorting:
The value of the sort parameter is a comma-separated list of sort keys. Supported sort keys include time, source, resource_type, resource_name, and event_type.
Each sort key may also include a direction. Supported directions are :asc for ascending and :desc for descending. The service will use :asc for every key that does not include a direction.
For example, to sort the list from most recently created to oldest:
GET /v1/events?sort=time:desc`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := auth.GetToken(keystoneDriver)
		if !token.Require("event:list") {
			return errors.New("You are not authorised to list events")
		}

		eventSlice, total, err := hermes.GetEvents(&data.Filter{}, &token.Context, keystoneDriver, storageDriver)
		if err != nil {
			return err
		}

		fmt.Printf("Total hits: %d\n", total)
		headers := []string{"Source", "Event ID", "Event Type", "Event Time", "Resource Name", "Resource Type", "User Name"}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)
		table.SetBorder(true)
		for _, ev := range eventSlice {
			dataRow := []string{ev.Source, ev.ID, ev.Type, ev.Time, ev.ResourceName, ev.ResourceType, ev.Initiator.UserName}
			table.Append(dataRow)
		}
		table.Render()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("source", "s", "", "Selects all events with this source.")
	listCmd.Flags().StringP("resource_type", "r", "", "Selects all events with this resource type.")
	listCmd.Flags().StringP("resource_name", "n", "", "Selects all events with this resource name.")
	listCmd.Flags().StringP("user_name", "u", "", "Selects all events with this user name.")
	listCmd.Flags().StringP("event_type", "e", "", "Selects all events with this event type.")
	listCmd.Flags().StringP("time", "t", "", "Date filter to select all events with event_time matching the specified criteria. See above for more detail.")
	listCmd.Flags().Int32P("offset", "o", 0, "The starting index within the total list of the events that you would like to retrieve..")
	listCmd.Flags().Int32P("limit", "l", 10, "The maximum number of records to return (up to 100). The default limit is 10.")
	listCmd.Flags().String("sort", "", "Determines the sorted order of the returned list. See above for more detail.")

}
