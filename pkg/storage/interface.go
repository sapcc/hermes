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

package storage

// Status contains Prometheus status strings
// TODO: Determine if we want a similar setup for Elasticsearch.
type Status string

const (
	// StatusSuccess means success
	StatusSuccess Status = "success"
	// StatusError means error
	StatusError = "error"
)

// ErrorType enumerates different Prometheus error types
type ErrorType string

const (
	// ErrorNone means no error
	ErrorNone ErrorType = ""
	// ErrorTimeout means that a timeout occurred while processing the request
	ErrorTimeout = "timeout"
	// ErrorCanceled means that the query was cancelled (to protect the service from malicious requests)
	ErrorCanceled = "canceled"
	// ErrorExec means unspecified error happened during query execution
	ErrorExec = "execution"
	// ErrorBadData means the API parameters where invalid
	ErrorBadData = "bad_data"
	// ErrorInternal means some unspecified internal error happened
	ErrorInternal = "internal"
)

// Response encapsulates a generic response of a Prometheus API
type Response struct {
	Status    Status        `json:"status"`
	Data      []interface{} `json:"data,omitempty"`
	ErrorType ErrorType     `json:"errorType,omitempty"`
	Error     string        `json:"error,omitempty"`
}

// Storage is an interface that wraps the underlying event storage mechanism.
// Because it is an interface, the real implementation can be mocked away in unit tests.
type Storage interface {
	/********** requests to ElasticSearch **********/
	GetEvents(filter *EventFilter, tenantID string) ([]*EventDetail, int, error)
	GetEvent(eventID string, tenantID string) (*EventDetail, error)
	GetAttributes(filter *AttributeFilter, tenantID string) ([]string, error)
	MaxLimit() uint
}

// FieldOrder maps the sort Fieldname and Order
type FieldOrder struct {
	Fieldname string
	Order     string //asc or desc
}

// EventFilter is similar to hermes.EventFilter, but using IDs instead of names
type EventFilter struct {
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

// AttributeFilter contains parameters for filtering by attributes
type AttributeFilter struct {
	QueryName string
	MaxDepth  uint
	Limit     uint
}

// Thanks to the tool at https://mholt.github.io/json-to-go/

//  The JSON annotations are for parsing the result from ElasticSearch
type eventListWithTotal struct {
	Total  int           `json:"total"`
	Events []EventDetail `json:"events"`
}

// Resource contains attributes describing a (OpenStack-) Resource
type Resource struct {
	TypeURI   string `json:"typeURI"`
	Name      string `json:"name,omitempty"`
	Domain    string `json:"domain,omitempty"`
	ID        string `json:"id"`
	Addresses []struct {
		URL  string `json:"url"`
		Name string `json:"name,omitempty"`
	} `json:"addresses,omitempty"`
	Host *struct {
		ID       string `json:"id,omitempty"`
		Address  string `json:"address,omitempty"`
		Agent    string `json:"agent,omitempty"`
		Platform string `json:"platform,omitempty"`
	} `json:"host,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// project_id and domain_id are OpenStack extensions (introduced by Keystone and keystone(audit)middleware)
	ProjectID string `json:"project_id,omitempty"`
	DomainID  string `json:"domain_id,omitempty"`
}

// Attachment contains self-describing extensions to the event
type Attachment struct {
	// Note: name is optional in CADF spec. to permit unnamed attachments
	Name string `json:"name,omitempty"`
	// this is messed-up in the spec.: the schema and examples says contentType. But the text often refers to typeURI.
	// Using typeURI would surely be more consistent. OpenStack uses typeURI, IBM supports both
	// (but forgot the name property)
	TypeURI string `json:"typeURI"`
	// Content contains the payload of the attachment. In theory this means any type.
	// In practise we have to decide because otherwise ES does based one first value
	// An interface allows arrays of json content. This should be json in the content.
	Content interface{} `json:"content"`
}

// EventDetail contains the CADF event according to CADF spec, section 6.6.1 Event (data)
// Extensions: requestPath (OpenStack, IBM), initiator.project_id/domain_id
// Omissions: everything that we do not use or not expose to API users
//  The JSON annotations are for parsing the result from ElasticSearch AND for generating the Hermes API response
type EventDetail struct {
	TypeURI   string `json:"typeURI"`
	ID        string `json:"id"`
	EventTime string `json:"eventTime"`
	Action    string `json:"action"`
	EventType string `json:"eventType"`
	Outcome   string `json:"outcome"`
	Reason    struct {
		ReasonType string `json:"reasonType"`
		ReasonCode string `json:"reasonCode"`
	} `json:"reason,omitempty"`
	Initiator   Resource     `json:"initiator"`
	Target      Resource     `json:"target"`
	Observer    Resource     `json:"observer"`
	Attachments []Attachment `json:"attachments,omitempty"`
	// requestPath is an extension of OpenStack's pycadf which is supported by IBM as well
	RequestPath string `json:"requestPath,omitempty"`
}

//AttributeValueList is used for holding unique attributes
type AttributeValueList []AttributeValue

//AttributeValue contains the attribute, and the number of hits.
type AttributeValue struct {
	Value string `json:"value"`
	count int64  `json:"-"`  // Removing export due to desire to not include it in JSON return
}
