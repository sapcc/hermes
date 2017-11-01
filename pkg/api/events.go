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

package api

import (
	"net/http"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sapcc/hermes/pkg/hermes"
	"github.com/sapcc/hermes/pkg/util"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// EventList is the model for JSON returned by the ListEvents API call
type EventList struct {
	NextURL string              `json:"next,omitempty"`
	PrevURL string              `json:"previous,omitempty"`
	Events  []*hermes.ListEvent `json:"events"`
	Total   int                 `json:"total"`
}

//ListEvents handles GET /v1/events.
func (p *v1Provider) ListEvents(res http.ResponseWriter, req *http.Request) {
	util.LogDebug("* api.ListEvents: Check token")
	token := p.CheckToken(req)
	if !token.Require(res, "event:list") {
		return
	}

	// Figure out the data.Filter to use, based on the request parameters

	// First off, parse the integers for offset & limit
	offset, _ := strconv.ParseUint(req.FormValue("offset"), 10, 32)
	limit, _ := strconv.ParseUint(req.FormValue("limit"), 10, 32)

	// Parse the sort querystring
	//slice of a struct, key and direction.

	sortSpec := []hermes.FieldOrder{}
	validSortTopics := map[string]bool{"time": true, "source": true, "resource_type": true, "resource_name": true, "event_type": true}
	validSortDirection := map[string]bool{"asc": true, "desc": true}
	sortParam := req.FormValue("sort")

	if sortParam != "" {
		for _, sortElement := range strings.Split(sortParam, ",") {
			keyVal := strings.SplitN(sortElement, ":", 2)
			//`time`, `source`, `resource_type`, `resource_name`, and `event_type`.
			sortfield := keyVal[0]
			if !validSortTopics[sortfield] {
				err := fmt.Errorf("Not a valid topic: %s, Valid topics: %v", sortfield, reflect.ValueOf(validSortTopics).MapKeys())
				http.Error(res, err.Error(), 400)
				return
			}

			defsortorder := "asc"
			if len(keyVal) == 2 {
				sortDirection := keyVal[1]
				if !validSortDirection[sortDirection] {
					err := fmt.Errorf("sort direction %s is invalid, must be asc or desc", sortDirection)
					http.Error(res, err.Error(), 400)
					return
				}
				defsortorder = sortDirection
			}

			s := hermes.FieldOrder{Fieldname: sortfield, Order: defsortorder}
			sortSpec = append(sortSpec, s)

		}
	}

	// Next, parse the elements of the time range filter
	timeRange := make(map[string]string)
	validOperators := map[string]bool{"lt": true, "lte": true, "gt": true, "gte": true}
	timeParam := req.FormValue("time")
	if timeParam != "" {
		for _, timeElement := range strings.Split(timeParam, ",") {
			keyVal := strings.SplitN(timeElement, ":", 2)
			operator := keyVal[0]
			if !validOperators[operator] {
				err := fmt.Errorf("time operator %s is not valid. Must be lt, lte, gt or gte", operator)
				http.Error(res, err.Error(), 400)
				return
			}
			_, exists := timeRange[operator]
			if exists {
				err := fmt.Errorf("Time operator %s can only occur once", operator)
				http.Error(res, err.Error(), 400)
				return
			}
			if len(keyVal) != 2 {
				err := fmt.Errorf("Time operator %s missing :<timestamp>", operator)
				http.Error(res, err.Error(), 400)
				return
			}
			validTimeFormats := []string{time.RFC3339, "2006-01-02T15:04:05-0700", "2006-01-02T15:04:05"}
			var isValidTimeFormat bool
			timeStr := keyVal[1]
			for _, timeFormat := range validTimeFormats {
				_, err := time.Parse(timeFormat, timeStr)
				if err != nil {
					isValidTimeFormat = true
					break
				}
			}
			if !isValidTimeFormat {
				err := fmt.Errorf("Invalid time format: %s", timeStr)
				http.Error(res, err.Error(), 400)
				return
			}
			timeRange[operator] = timeStr
		}
	}

	util.LogDebug("api.ListEvents: Create filter")
	filter := hermes.Filter{
		// TODO: remove append of deprecated parameters
		ObserverType: req.FormValue("observer_type") + req.FormValue("source"),
		TargetType:   req.FormValue("target_type") + req.FormValue("resource_type"),
		TargetID:     req.FormValue("target_id"),
		OriginatorID: req.FormValue("originator_id") + req.FormValue("user_name"),
		Action:       req.FormValue("action") + req.FormValue("event_type"),
		Time:         timeRange,
		Offset:       uint(offset),
		Limit:        uint(limit),
		Sort:         sortSpec,
	}

	util.LogDebug("api.ListEvents: call hermes.GetEvents()")
	tenantID, err := getTenantID(token, req, res)
	if err != nil {
		return
	}
	events, total, err := hermes.GetEvents(&filter, tenantID, p.keystone, p.storage)
	if ReturnError(res, err) {
		util.LogError("api.ListEvents: error %s", err)
		return
	}

	eventList := EventList{Events: events, Total: total}

	// What protocol to use for PrevURL and NextURL?
	protocol := getProtocol(req)
	// Do we need a NextURL?
	if int(filter.Offset+filter.Limit) < total {
		req.Form.Set("offset", strconv.FormatUint(uint64(filter.Offset+filter.Limit), 10))
		eventList.NextURL = fmt.Sprintf("%s://%s%s?%s", protocol, req.Host, req.URL.Path, req.Form.Encode())
	}
	// Do we need a PrevURL?
	if int(filter.Offset-filter.Limit) >= 0 {
		req.Form.Set("offset", strconv.FormatUint(uint64(filter.Offset-filter.Limit), 10))
		eventList.PrevURL = fmt.Sprintf("%s://%s%s?%s", protocol, req.Host, req.URL.Path, req.Form.Encode())
	}

	ReturnJSON(res, 200, eventList)
}
func getProtocol(req *http.Request) string {
	protocol := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	return protocol
}

//GetEvent handles GET /v1/events/:event_id.
func (p *v1Provider) GetEventDetails(res http.ResponseWriter, req *http.Request) {
	token := p.CheckToken(req)
	if !token.Require(res, "event:show") {
		return
	}
	eventID := mux.Vars(req)["event_id"]
	tenantID, err := getTenantID(token, req, res)
	if err != nil {
		return
	}

	event, err := hermes.GetEvent(eventID, tenantID, p.keystone, p.storage)

	if ReturnError(res, err) {
		return
	}
	if event == nil {
		err := fmt.Errorf("Event %s could not be found in tenant %s", eventID, tenantID)
		http.Error(res, err.Error(), 404)
		return
	}
	ReturnJSON(res, 200, event)
}

//GetAttributes handles GET /v1/attributes/:attribute_name
func (p *v1Provider) GetAttributes(res http.ResponseWriter, req *http.Request) {
	token := p.CheckToken(req)
	if !token.Require(res, "event:show") {
		return
	}
	queryName := mux.Vars(req)["attribute_name"]
	if queryName == "" {
		fmt.Errorf("QueryName not found")
	}
	tenantID, err := getTenantID(token, req, res)
	if err != nil {
		return
	}

	attribute, err := hermes.GetAttributes(queryName, tenantID, p.storage)

	if ReturnError(res, err) {
		return
	}
	if attribute == nil {
		err := fmt.Errorf("Attribute %s could not be found in tenant %s", attribute, tenantID)
		http.Error(res, err.Error(), 404)
		return
	}
	ReturnJSON(res, 200, attribute)
}

func getTenantID(token *Token, r *http.Request, w http.ResponseWriter) (string, error) {
	// Get tenant id from token
	tenantID := token.context.Auth["tenant_id"]
	if tenantID == "" {
		tenantID = token.context.Auth["domain_id"]
	}
	// Tenant id can be overriden with a query parameter
	projectID := r.FormValue("project_id")
	domainID := r.FormValue("domain_id")
	if projectID != "" {
		tenantID = projectID
	}
	if domainID != "" {
		if projectID != "" {
			err := errors.New("domain_id and project_id cannot both be specified")
			http.Error(w, err.Error(), 400)
			return "", err
		}
		tenantID = domainID
	}
	return tenantID, nil
}
