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
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/sapcc/hermes/pkg/identity"
	"github.com/sapcc/hermes/pkg/storage"
	"github.com/sapcc/hermes/pkg/util"
)

//ResourceRef is an embedded struct for ListEvents (eg. Initiator, Target, Observer)
type ResourceRef struct {
	TypeURI string `json:"typeURI"`
	ID      string `json:"id"`
}

// ListEvent contains high-level data about an event, intended as a list item
//  The JSON annotations here are for the JSON to be returned by the API
type ListEvent struct {
	// TODO: Remove these deprecations
	SourceDeprecated       string `json:"source"`                  // observer.typeURI
	IDDeprecated           string `json:"event_id"`                // id
	TypeDeprecated         string `json:"event_type"`              // action
	TimeDeprecated         string `json:"event_time"`              // eventType
	ResourceIDDeprecated   string `json:"resource_id"`             // target.id
	ResourceTypeDeprecated string `json:"resource_type"`           // target.typeURI
	ResourceNameDeprecated string `json:"resource_name,omitempty"` // drop
	// NEW:
	ID        string      `json:"id"`
	Time      string      `json:"eventTime"`
	Action    string      `json:"action"`
	Outcome   string      `json:"outcome"`
	Initiator ResourceRef `json:"initiator"`
	Target    ResourceRef `json:"target"`
	Observer  ResourceRef `json:"observer"`
}

// FieldOrder maps the sort Fieldname and Order
type FieldOrder struct {
	Fieldname string
	Order     string //asc or desc
}

// Filter maps to the filtering/paging/sorting allowed by the API
type Filter struct {
	ObserverType  string
	TargetType    string
	TargetID      string
	InitiatorID   string
	InitiatorType string
	Action        string
	Outcome       string
	Time          map[string]string
	Offset        uint
	Limit         uint
	Sort          []FieldOrder
}

// GetEvents returns a list of matching events (with filtering)
func GetEvents(filter *Filter, tenantID string, keystoneDriver identity.Identity, eventStore storage.Storage) ([]*ListEvent, int, error) {
	storageFilter, err := storageFilter(filter, keystoneDriver, eventStore)
	if err != nil {
		return nil, 0, err
	}
	util.LogDebug("hermes.GetEvents: tenant id is %s", tenantID)
	eventDetails, total, err := eventStore.GetEvents(storageFilter, tenantID)
	if err != nil {
		return nil, 0, err
	}
	events, err := eventsList(eventDetails, keystoneDriver)
	if err != nil {
		return nil, 0, err
	}
	return events, total, err
}

func storageFilter(filter *Filter, keystoneDriver identity.Identity, eventStore storage.Storage) (*storage.Filter, error) {
	// As per the documentation, the default limit is 10
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	if filter.Offset+filter.Limit > eventStore.MaxLimit() {
		return nil, fmt.Errorf("offset %d plus limit %d exceeds the maximum of %d",
			filter.Offset, filter.Limit, eventStore.MaxLimit())
	}

	storageFieldOrder := []storage.FieldOrder{}
	err := copier.Copy(&storageFieldOrder, &filter.Sort)
	if err != nil {
		panic("Could not copy storage field order.")
	}
	storageFilter := storage.Filter{
		ObserverType:  filter.ObserverType,
		InitiatorID:   filter.InitiatorID,
		InitiatorType: filter.InitiatorType,
		TargetType:    filter.TargetType,
		TargetID:      filter.TargetID,
		Action:        filter.Action,
		Outcome:       filter.Outcome,
		Time:          filter.Time,
		Offset:        filter.Offset,
		Limit:         filter.Limit,
		Sort:          storageFieldOrder,
	}
	return &storageFilter, nil
}

// Construct ListEvents - Optionally (default off) add the names for IDs in the events
func eventsList(eventDetails []*storage.EventDetail, keystoneDriver identity.Identity) ([]*ListEvent, error) {
	var events []*ListEvent
	for _, storageEvent := range eventDetails {
		event := ListEvent{
			// TODO: remove old attribute names
			SourceDeprecated:       storageEvent.Observer.TypeURI,
			IDDeprecated:           storageEvent.ID,
			TypeDeprecated:         storageEvent.Action,
			TimeDeprecated:         storageEvent.EventTime,
			ResourceIDDeprecated:   storageEvent.Target.ID,
			ResourceTypeDeprecated: storageEvent.Target.TypeURI,
			// new attributes
			Initiator: ResourceRef{
				TypeURI: storageEvent.Initiator.TypeURI,
				ID:      storageEvent.Initiator.ID,
			},
			Target: ResourceRef{
				TypeURI: storageEvent.Target.TypeURI,
				ID:      storageEvent.Target.ID,
			},
			ID:      storageEvent.ID,
			Action:  storageEvent.Action,
			Outcome: storageEvent.Outcome,
			Time:    storageEvent.EventTime,
			Observer: ResourceRef{
				TypeURI: storageEvent.Observer.TypeURI,
				ID:      storageEvent.Observer.ID,
			},
		}
		err := copier.Copy(&event.Initiator, &storageEvent.Initiator)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}
	return events, nil
}

// GetEvent returns the CADF detail for event with the specified ID
func GetEvent(eventID string, tenantID string, keystoneDriver identity.Identity, eventStore storage.Storage) (*storage.EventDetail, error) {
	event, err := eventStore.GetEvent(eventID, tenantID)

	return event, err
}

//GetAttributes No Logic here, but handles mock implementation for eventStore
func GetAttributes(queryName string, tenantID string, eventStore storage.Storage) ([]string, error) {
	attribute, err := eventStore.GetAttributes(queryName, tenantID)

	return attribute, err
}
