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
	// ErrorTimeout means that a timeout occured while processing the request
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

// AttributeFilter
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

// EventDetail contains the CADF payload, enhanced with names for IDs
//  The JSON annotations are for parsing the result from ElasticSearch AND for generating the Hermes API response
type EventDetail struct {
	TypeURI   string `json:"typeURI"`
	ID        string `json:"id"`
	EventTime string `json:"eventTime"`
	Action    string `json:"action"`
	EventType string `json:"eventType"`
	Outcome   string `json:"outcome"`
	Initiator struct {
		TypeURI   string `json:"typeURI"`
		ID        string `json:"id"`
		DomainID  string `json:"domain_id,omitempty"`
		ProjectID string `json:"project_id,omitempty"`
		Host      struct {
			Agent   string `json:"agent"`
			Address string `json:"address"`
		} `json:"host,omitempty"`
	} `json:"initiator"`
	Target struct {
		TypeURI     string `json:"typeURI"`
		ID          string `json:"id"`
		Name        string `json:"name,omitempty"`
		ProjectID   string `json:"project_id,omitempty"`
		DomainID    string `json:"domain_id,omitempty"`
		Attachments []struct {
			Content     string `json:"content"`
			ContentType string `json:"contentType"`
			Name        string `json:"name,omitempty"`
		} `json:"attachments,omitempty"`
	} `json:"target"`
	Observer struct {
		TypeURI string `json:"typeURI"`
		ID      string `json:"id"`
	} `json:"observer"`
	// TODO: attachments for additional parameters
}

//AttributeValueList is used for holding unique attributes
type AttributeValueList []AttributeValue

//AttributeValue contains the attribute, and the number of hits.
type AttributeValue struct {
	Value string `json:"value"`
	count int64  `json:"count"` // Removing export due to desire to not include it in JSON return
}
