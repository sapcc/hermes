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
	"errors"
	"fmt"
	"github.com/databus23/goslo.policy"
	"github.com/prometheus/common/log"
	"github.com/sapcc/hermes/pkg/data"
	"github.com/sapcc/hermes/pkg/keystone"
	"github.com/sapcc/hermes/pkg/storage"
)

// GetEvents returns a list of matching events (with filtering)
func GetEvents(filter *data.Filter, context *policy.Context, eventStore storage.Interface) ([]*data.Event, int, error) {
	if context == nil {
		return nil, 0, errors.New("GetEvent() called with no policy context")
	}

	// As per the documentation, the default limit is 10
	if filter.Limit == 0 {
		filter.Limit = 10
	}

	events, total, err := eventStore.GetEvents(*filter, getTenantId(context))

	// Now add the names for IDs in the events
	keystoneSvc := keystone.ConfiguredDriver()
	for _, event := range events {
		if err == nil && len(event.Initiator.DomainID) != 0 {
			event.Initiator.DomainName, err = keystoneSvc.DomainName(event.Initiator.DomainID)
			if err != nil {
				log.Errorf("Error looking up domain name for domain '%s'", event.Initiator.DomainID)
			}
		}
		if err == nil && len(event.Initiator.ProjectID) != 0 {
			event.Initiator.ProjectName, err = keystoneSvc.ProjectName(event.Initiator.ProjectID)
			if err != nil {
				log.Errorf("Error looking up project name for project '%s'", event.Initiator.ProjectID)
			}
		}
		if err == nil && len(event.Initiator.UserID) != 0 {
			event.Initiator.UserName, err = keystoneSvc.UserName(event.Initiator.UserID)
			if err != nil {
				log.Errorf("Error looking up user name for user '%s'", event.Initiator.UserID)
			}
		}

		// Depending on the type of the target, we need to look up the name in different services
		if err == nil {
			switch event.ResourceType {
			case "data/security/project":
				event.ResourceName, err = keystoneSvc.ProjectName(event.ResourceId)
			default:
				log.Warn(fmt.Sprintf("Unhandled payload type \"%s\", cannot look up name.",
					event.ResourceType))
			}
		}

		if err != nil {
			break
		}
	}

	return events, total, err
}

// GetEvent returns the CADF detail for event with the specified ID
func GetEvent(eventID string, context *policy.Context, eventStore storage.Interface) (*data.EventDetail, error) {
	if context == nil {
		return nil, errors.New("GetEvent() called with no policy context")
	}
	event, err := eventStore.GetEvent(eventID, getTenantId(context))
	// Now add the names for IDs in the event
	keystoneSvc := keystone.ConfiguredDriver()
	if err == nil && event.Payload.Initiator.DomainID != "" {
		event.Payload.Initiator.DomainName, err = keystoneSvc.DomainName(event.Payload.Initiator.DomainID)
	}
	if err == nil && event.Payload.Initiator.ProjectID != "" {
		event.Payload.Initiator.ProjectName, err = keystoneSvc.ProjectName(event.Payload.Initiator.ProjectID)
	}
	if err == nil && event.Payload.Initiator.UserID != "" {
		event.Payload.Initiator.UserName, err = keystoneSvc.UserName(event.Payload.Initiator.UserID)
	}

	// Depending on the type of the target, we need to look up the name in different services
	if err == nil {
		switch event.Payload.Target.TypeURI {
		case "data/security/project":
			event.Payload.Target.Name, err = keystoneSvc.ProjectName(event.Payload.Target.ID)
		default:
			log.Warn(fmt.Sprintf("Unhandled payload type \"%s\", cannot look up name.",
				event.Payload.Target.TypeURI))
		}
	}
	return &event, err
}

func getTenantId(context *policy.Context) string {
	id, project := context.Auth["project_id"]
	if project {
		return id
	}
	return context.Auth["domain_id"]
}
