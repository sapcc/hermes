package storage

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/sapcc/hermes/pkg/data"
	"strings"
)

type mock struct{}

func Mock() Interface {
	return mock{}
}

type eventsList struct {
	total  int
	events []data.EventDetail
}

func (m mock) GetEvents(filter data.Filter, tenant_id string) ([]*data.Event, int, error) {
	var detailedEvents eventListWithTotal
	json.Unmarshal(mockEvents, &detailedEvents)

	var events []*data.Event

	for _, de := range detailedEvents.Events {
		p := de.Payload
		ev := data.Event{
			Source:       strings.SplitN(de.EventType, ".", 2)[0],
			ID:           p.ID,
			Type:         de.EventType,
			Time:         p.EventTime,
			ResourceId:   de.Payload.Target.ID,
			ResourceType: de.Payload.Target.TypeURI,
		}
		err := copier.Copy(&ev.Initiator, &de.Payload.Initiator)
		if err != nil {
			return nil, 0, nil
		}
		events = append(events, &ev)
	}

	return events, detailedEvents.Total, nil
}

func (m mock) GetEvent(eventId string, tenant_id string) (data.EventDetail, error) {
	var parsedEvent Event
	json.Unmarshal(mockEvent, &parsedEvent)
	event := data.EventDetail{}
	err := copier.Copy(&event, &parsedEvent)
	return event, err
}

var mockEvent = []byte(`
{
	"publisher_id": "identity.keystone-2031324599-gujvn",
	"event_type": "identity.project.deleted",
	"payload": {
		"observer": {
			"typeURI": "service/security",
			"id": "493f1d6d-af50-5a4b-813b-488ecdfb1010"
		},
		"resource_info": "b3b70c8271a845709f9a03030e705da7",
		"typeURI": "http://schemas.dmtf.org/cloud/audit/1.0/event",
		"initiator": {
			"typeURI": "service/security/account/user",
			"project_id": "6a030751147a45c0863c3b5bde32c744",
			"user_id": "eb5cd8f904b06e8b2a6eb86c8b04c08e6efb89b92da77905cc8c475f30b0b812",
			"host": {
				"agent": "python-keystoneclient",
				"address": "100.65.0.11"
			},
			"id": "4a70d16f08b05d038c1e5ee7a5ee554e"
		},
		"eventTime": "2017-05-02T12:02:46.726056+0000",
		"action": "deleted.project",
		"eventType": "activity",
		"id": "d5eed458-6666-58ec-ad06-8d3cf6bafca1",
		"outcome": "success",
		"target": {
			"typeURI": "data/security/project",
			"id": "b3b70c8271a845709f9a03030e705da7"
		}
	},
	"message_id": "5a32c2f3-2996-4f46-819c-6197cf06037e",
	"priority": "info",
	"timestamp": "2017-05-02 12:02:46.726619"
}
`)

var mockEvents = []byte(`
{
	"total": 24,
	"events": [{
			"publisher_id": "identity.keystone-2031324599-gujvn",
			"event_type": "identity.project.deleted",
			"payload": {
				"observer": {
					"typeURI": "service/security",
					"id": "493f1d6d-af50-5a4b-813b-488ecdfb1010"
				},
				"resource_info": "b3b70c8271a845709f9a03030e705da7",
				"typeURI": "http://schemas.dmtf.org/cloud/audit/1.0/event",
				"initiator": {
					"typeURI": "service/security/account/user",
					"project_id": "ae63ddf2076d4342a56eb049e37a7621",
					"user_id": "eb5cd8f904b06e8b2a6eb86c8b04c08e6efb89b92da77905cc8c475f30b0b812",
					"host": {
						"agent": "python-keystoneclient",
						"address": "100.65.0.11"
					},
					"id": "4a70d16f08b05d038c1e5ee7a5ee554e"
				},
				"eventTime": "2017-05-02T12:02:46.726056+0000",
				"action": "deleted.project",
				"eventType": "activity",
				"id": "d5eed458-6666-58ec-ad06-8d3cf6bafca1",
				"outcome": "success",
				"target": {
					"typeURI": "data/security/project",
					"id": "b3b70c8271a845709f9a03030e705da7"
				}
			},
			"message_id": "5a32c2f3-2996-4f46-819c-6197cf06037e",
			"priority": "info",
			"timestamp": "2017-05-02 12:02:46.726619"
		}, {
			"publisher_id": "identity.keystone-2031324599-gujvn",
			"event_type": "identity.project.deleted",
			"payload": {
				"observer": {
					"typeURI": "service/security",
					"id": "a66f7b00-b52d-51a1-b370-4e129bd534e2"
				},
				"resource_info": "b3b70c8271a845709f9a03030e705da7",
				"typeURI": "http://schemas.dmtf.org/cloud/audit/1.0/event",
				"initiator": {
					"typeURI": "service/security/account/user",
					"project_id": "ae63ddf2076d4342a56eb049e37a7621",
					"user_id": "eb5cd8f904b06e8b2a6eb86c8b04c08e6efb89b92da77905cc8c475f30b0b812",
					"host": {
						"agent": "python-keystoneclient",
						"address": "100.64.0.4"
					},
					"id": "4a70d16f08b05d038c1e5ee7a5ee554e"
				},
				"eventTime": "2017-05-02T11:45:49.982112+0000",
				"action": "deleted.project",
				"eventType": "activity",
				"id": "095056c9-4cbb-5200-af70-0977dbcf5000",
				"outcome": "success",
				"target": {
					"typeURI": "data/security/project",
					"id": "b3b70c8271a845709f9a03030e705da7"
				}
			},
			"message_id": "c3c61a95-54f9-44d0-9986-9571258646cd",
			"priority": "info",
			"timestamp": "2017-05-02 11:45:49.982909"
		}, {
			"publisher_id": "identity.keystone-2031324599-gujvn",
			"event_type": "identity.project.deleted",
			"payload": {
				"observer": {
					"typeURI": "service/security",
					"id": "15276db2-9b34-528c-b72a-7eca6995bf58"
				},
				"resource_info": "b3b70c8271a845709f9a03030e705da7",
				"typeURI": "http://schemas.dmtf.org/cloud/audit/1.0/event",
				"initiator": {
					"typeURI": "service/security/account/user",
					"project_id": "ae63ddf2076d4342a56eb049e37a7621",
					"user_id": "eb5cd8f904b06e8b2a6eb86c8b04c08e6efb89b92da77905cc8c475f30b0b812",
					"host": {
						"agent": "python-keystoneclient",
						"address": "100.64.0.4"
					},
					"id": "4a70d16f08b05d038c1e5ee7a5ee554e"
				},
				"eventTime": "2017-05-02T11:45:44.755215+0000",
				"action": "deleted.project",
				"eventType": "activity",
				"id": "dbd72ad7-61b4-5dab-b9ed-26068a187c7a",
				"outcome": "success",
				"target": {
					"typeURI": "data/security/project",
					"id": "b3b70c8271a845709f9a03030e705da7"
				}
			},
			"message_id": "0cd52307-f09f-453f-bf1b-027b2f907e94",
			"priority": "info",
			"timestamp": "2017-05-02 11:45:44.756160"
		}
	]
}
`)
