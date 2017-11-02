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
	ResourceIdDeprecated   string `json:"resource_id"`             // target.id
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
	// TODO: double-check if we really want to do this (think of all the different resource types out there)
	// IMHO (jobrs) resolving IDs into names can be done on the UI and on request of the user
	//// Translate hermes.Filter to storage.Filter by filling in IDs for names
	//if filter.ResourceName != "" {
	//	// TODO: make sure there is a resource type, then look up the corresponding name
	//	//storageFilter.TargetID = resourceId
	//}
	//if filter.InitiatorID != "" {
	//	util.LogDebug("Filtering on InitiatorID: %s", filter.InitiatorID)
	//	//userId, err := keystoneDriver.InitiatorID(filter.InitiatorID)
	//	//if err != nil {
	//	//	util.LogError("Could not find user ID &s for name %s", userId, filter.InitiatorID)
	//	//}
	//	storageFilter.InitiatorID = filter.InitiatorID
	//}
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
			ResourceIdDeprecated:   storageEvent.Target.ID,
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

		// TODO: adapt and reactivate if needed
		//if viper.GetBool("hermes.enrich_keystone_events") {
		//	nameMap := namesForIds(keystoneDriver, map[string]string{
		//		"init_user_domain":  event.Initiator.DomainID,
		//		"init_user_project": event.Initiator.ProjectID,
		//		"init_user":         event.Initiator.UserID,
		//		"target":            event.TargetID,
		//	}, event.TargetType)
		//
		//	//event.Initiator.DomainName = nameMap["init_user_domain"]
		//	//event.Initiator.ProjectName = nameMap["init_user_project"]
		//	//event.Initiator.InitiatorID = nameMap["init_user"]
		//	event.ResourceName = nameMap["target"]
		//}
		events = append(events, &event)
	}
	return events, nil
}

// GetEvent returns the CADF detail for event with the specified ID
func GetEvent(eventID string, tenantID string, keystoneDriver identity.Identity, eventStore storage.Storage) (*storage.EventDetail, error) {
	event, err := eventStore.GetEvent(eventID, tenantID)

	/* TODO: think about whether this makes sense. In CADF, arbitrary resources are referenced by typeURI and ID.
	Should we attempt to resolve these IDs into names or leave this up to the UI layer? The UI layer should have that
	functionality and corresponding caching mechanism already. Also names are not unique and can change, so they cannot
	be used safely for anything but presentation.
	if viper.GetBool("hermes.enrich_keystone_events") {
		if event != nil {
			nameMap := namesForIds(keystoneDriver, map[string]string{
				"init_user_domain":  event.Payload.Initiator.DomainID,
				"init_user_project": event.Payload.Initiator.ProjectID,
				"init_user":         event.Payload.Initiator.UserID,
				"target":            event.Payload.Target.ID,
				"project":           event.Payload.Project,
				"user":              event.Payload.User,
				"group":             event.Payload.Group,
				"role":              event.Payload.Role,
			}, event.Payload.Target.TypeURI)

			event.Initiator.DomainName = nameMap["init_user_domain"]
			event.Initiator.ProjectName = nameMap["init_user_project"]
			event.Initiator.InitiatorID = nameMap["init_user"]
			event.Target.Name = nameMap["target"]
			event.ProjectName = nameMap["project"]
			event.InitiatorID = nameMap["user"]
			event.GroupName = nameMap["group"]
			event.RoleName = nameMap["role"]
		}
	}
	*/
	return event, err
}

//GetAttributes No Logic here, but handles mock implementation for eventStore
func GetAttributes(queryName string, tenantID string, eventStore storage.Storage) ([]string, error) {
	attribute, err := eventStore.GetAttributes(queryName, tenantID)

	return attribute, err
}

// TODO: remove or extend for all those Nova, Neutron, ... resource types
/*
func namesForIds(keystoneDriver identity.Identity, idMap map[string]string, targetType string) map[string]string {
	nameMap := map[string]string{}
	var err error

	// Now add the names for IDs in the event to the nameMap
	iUserDomainID := idMap["init_user_domain"]
	if iUserDomainID != "" {
		nameMap["init_user_domain"], err = keystoneDriver.DomainName(iUserDomainID)
		if err != nil {
			log.Printf("Error looking up domain name for domain '%s'", iUserDomainID)
		}
	}
	iUserProjectID := idMap["init_user_project"]
	if iUserProjectID != "" {
		nameMap["init_user_project"], err = keystoneDriver.ProjectName(iUserProjectID)
		if err != nil {
			log.Printf("Error looking up project name for project '%s'", iUserProjectID)
		}
	}
	iUserID := idMap["init_user"]
	if iUserID != "" {
		nameMap["init_user"], err = keystoneDriver.InitiatorID(iUserID)
		if err != nil {
			log.Printf("Error looking up user name for user '%s'", iUserID)
		}
	}
	projectID := idMap["project"]
	if projectID != "" {
		nameMap["project"], err = keystoneDriver.ProjectName(projectID)
		if err != nil {
			log.Printf("Error looking up project name for project '%s'", projectID)
		}
	}
	userID := idMap["user"]
	if userID != "" {
		nameMap["user"], err = keystoneDriver.InitiatorID(userID)
		if err != nil {
			log.Printf("Error looking up user name for user '%s'", userID)
		}
	}
	groupID := idMap["group"]
	if groupID != "" {
		nameMap["group"], err = keystoneDriver.GroupName(groupID)
		if err != nil {
			log.Printf("Error looking up user name for group '%s'", groupID)
		}
	}
	roleID := idMap["role"]
	if roleID != "" {
		nameMap["role"], err = keystoneDriver.RoleName(roleID)
		if err != nil {
			log.Printf("Error looking up user name for role '%s'", roleID)
		}
	}

	// Depending on the type of the target, we need to look up the name in different services
	switch targetType {
	case "data/security/project":
		nameMap["target"], err = keystoneDriver.ProjectName(idMap["target"])
	case "service/security/account/user":
	// doesn't work for users - a UUID is used for some reason, which can't be looked up
	//	nameMap["target"], err = keystoneDriver.InitiatorID(idMap["target"])
	default:
		log.Printf("Unhandled payload type \"%s\", cannot look up name.", targetType)
	}
	if err != nil {
		log.Printf("Error looking up name for %s '%s'", targetType, idMap["target"])
	}

	return nameMap
}
*/
