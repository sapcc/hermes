/*******************************************************************************
*
* Copyright 2017 SAP SE
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

package hermes

import (
	"github.com/sapcc/hermes/pkg/data"
	"github.com/sapcc/hermes/pkg/storage"
)

// GetEvents returns Event reports for all events (with filtering) or, if eventID is
// non-nil, the CADF detail for that event only.
func GetEvents(eventID *string, eventStore storage.Interface, filter *data.Filter) ([]*data.Event, error) {

	events := make(events)

	// TODO - call the backend storage driver to get the data
	result := make([]*data.Event, len(events))

	return result, nil
}

func makeEventFilter(tableWithEventID string, eventID *string) map[string]interface{} {
	fields := make(map[string]interface{})
	if eventID != nil {
		fields[tableWithEventID+".event_id"] = *eventID
	}
	return fields
}

type events map[string]*data.Event
