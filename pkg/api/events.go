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

	// Next, parse the elements of the time range filter
	timeRange := make(map[string] string)
	validOperators := map[string] bool{"lt": true, "lte": true, "gt": true, "gte": true}
	timeParam := req.FormValue("time")
	if timeParam != "" {
		for _, timeElement := range strings.Split(timeParam, ",") {
			keyVal := strings.SplitN(timeElement, ":", 2)
			operator := keyVal[0]
			if !validOperators[operator] {
				err := errors.New(fmt.Sprintf("Time operator %s is not valid. Must be lt, lte, gt or gte.", operator))
				http.Error(res, err.Error(), 400)
				return
			}
			_, exists := timeRange[operator]
			if exists {
				err := errors.New(fmt.Sprintf("Time operator %s can only occur once", operator))
				http.Error(res, err.Error(), 400)
				return
			}
			if len(keyVal) != 2 {
				err := errors.New(fmt.Sprintf("Time operator %s missing :<timestamp>", operator))
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
				err := errors.New(fmt.Sprintf("Invalid time format: %s", timeStr))
				http.Error(res, err.Error(), 400)
				return
			}
			timeRange[operator] = timeStr
		}
	}

	util.LogDebug("api.ListEvents: Create filter")
	filter := hermes.Filter{
		Source:       req.FormValue("source"),
		ResourceType: req.FormValue("resource_type"),
		ResourceName: req.FormValue("resource_name"),
		UserName:     req.FormValue("user_name"),
		EventType:    req.FormValue("event_type"),
		Time:         timeRange,
		Offset:       uint(offset),
		Limit:        uint(limit),
		Sort:         req.FormValue("sort"),
	}

	util.LogDebug("api.ListEvents: call hermes.GetEvents()")
	tenantId, err := getTenantId(token, req, res)
	if err != nil {
		return
	}
	events, total, err := hermes.GetEvents(&filter, tenantId, p.keystone, p.storage)
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
	tenantId, err := getTenantId(token, req, res)
	if err != nil {
		return
	}

	event, err := hermes.GetEvent(eventID, tenantId, p.keystone, p.storage)

	if ReturnError(res, err) {
		return
	}
	if event == nil {
		err := fmt.Errorf("Event %s could not be found in tenant %s", eventID, tenantId)
		http.Error(res, err.Error(), 404)
		return
	}
	ReturnJSON(res, 200, event)
}

func getTenantId(token *Token, r *http.Request, w http.ResponseWriter) (string, error) {
	// Get tenant id from token
	tenantId := token.context.Auth["tenant_id"]
	if tenantId=="" {
		tenantId = token.context.Auth["domain_id"]
	}
	// Tenant id can be overriden with a query parameter
	projectId := r.FormValue("project_id")
	domainId := r.FormValue("domain_id")
	if projectId != "" {
		tenantId = projectId
	}
	if domainId != "" {
		if projectId != "" {
			err := errors.New("domain_id and project_id cannot both be specified")
			http.Error(w, err.Error(), 400)
			return "", err
		}
		tenantId = domainId
	}
	return tenantId, nil
}
