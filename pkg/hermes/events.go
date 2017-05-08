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

// GetEvents returns a list of matching events (with filtering)
func GetEvents(eventStore storage.Interface, filter data.Filter) ([]data.Event, int, error) {
	events, total, error := eventStore.GetEvents(filter)
	return events, total, error
}

// GetEvent returns the CADF detail for event with the specified ID
func GetEvent(eventID string, eventStore storage.Interface) (data.EventDetail, error) {
	event, error := eventStore.GetEvent(eventID)
	return event, error
}
