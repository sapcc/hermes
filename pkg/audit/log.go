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

package audit

import (
	"encoding/json"
	"os"
	"time"

	"github.com/sapcc/go-bits/logg"
	"github.com/sapcc/go-bits/retry"
	"log"
)

func init() {
	log.SetOutput(os.Stdout)
	if os.Getenv("AUDIT_DEBUG") == "1" {
		logg.ShowDebug = true
	}
}

//Trail is a list of CADF formatted events with log level AUDIT. It has a separate interface
//from the rest of the logging to allow to withhold the logging until DB changes are committed.
type Trail struct {
	events []Event
}

// CADFEvent contains the CADF event according to CADF spec, section 6.6.1 Event (data)
// Extensions: requestPath (OpenStack, IBM), initiator.project_id/domain_id
// Omissions: everything that we do not use or not expose to API users
type Event struct {
	// CADF Event Schema
	TypeURI string `json:"typeURI"`

	// CADF generated event id
	ID string `json:"id"`

	// CADF generated timestamp
	EventTime string `json:"eventTime"`

	// Characterizes events: eg. activity
	EventType string `json:"eventType"`

	// CADF action mapping for GET call on an OpenStack REST API
	Action string `json:"action"`

	// Outcome of REST API call, eg. success/failure
	Outcome string `json:"outcome"`

	// Standard response for successful HTTP requests
	Reason Reason `json:"reason,omitempty"`

	// CADF component that contains the RESOURCE
	// that initiated, originated, or instigated the event's
	// ACTION, according to the OBSERVER
	Initiator Resource `json:"initiator"`

	// CADF component that contains the RESOURCE
	// against which the ACTION of a CADF Event
	// Record was performed, was attempted, or is
	// pending, according to the OBSERVER.
	Target Resource `json:"target"`

	// CADF component that contains the RESOURCE
	// that generates the CADF Event Record based on
	// its observation (directly or indirectly) of the Actual Event
	Observer Resource `json:"observer"`

	// Attachment contains self-describing extensions to the event
	Attachments []Attachment `json:"attachments,omitempty"`

	// Request path on the OpenStack service REST API call
	RequestPath string `json:"requestPath,omitempty"`
}

//Resource is a substructure of CADFEvent and contains attributes describing a (OpenStack-) resource.
type Resource struct {
	TypeURI   string `json:"typeURI"`
	Name      string `json:"name,omitempty"`
	Domain    string `json:"domain,omitempty"`
	ID        string `json:"id"`
	Addresses []struct {
		URL  string `json:"url"`
		Name string `json:"name,omitempty"`
	} `json:"addresses,omitempty"`
	Host        *Host        `json:"host,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	ProjectID   string       `json:"project_id,omitempty"`
	DomainID    string       `json:"domain_id,omitempty"`
}

//Attachment is a substructure of CADFEvent and contains self-describing extensions to the event.
type Attachment struct {
	Name    string      `json:"name,omitempty"`
	TypeURI string      `json:"typeURI"`
	Content interface{} `json:"content"`
}

//Reason is a substructure of CADFevent containing data for the event outcome's reason.
type Reason struct {
	ReasonType string `json:"reasonType"`
	ReasonCode string `json:"reasonCode"`
}

//Host is a substructure of Resource containing data for the event initiator's host.
type Host struct {
	ID       string `json:"id,omitempty"`
	Address  string `json:"address,omitempty"`
	Agent    string `json:"agent,omitempty"`
	Platform string `json:"platform,omitempty"`
}

//Add adds an event to the audit trail.
func (t *Trail) Add(event Event) {
	t.events = append(t.events, event)
}

//Commit sends the whole audit trail into the log. Call this after tx.Commit().
func (t *Trail) Commit(clusterID string, config Config) {
	if config.Enabled && len(t.events) != 0 {
		events := t.events //take a copy to pass into the goroutine
		go retry.ExponentialBackoff{
			Factor:      2,
			MaxInterval: 5 * time.Minute,
		}.RetryUntilSuccessful(func() error { return sendEvents(clusterID, config, events) })
	}

	for _, event := range t.events {
		msg, _ := json.Marshal(event)
		logg.Other("AUDIT", string(msg))
	}
	t.events = nil //do not log these lines again
}
