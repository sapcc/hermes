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
)

// ListEvent list for returning in the API
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
	offset, _ := strconv.ParseUint(req.FormValue("offset"), 10, 32)
	limit, _ := strconv.ParseUint(req.FormValue("limit"), 10, 8)

	util.LogDebug("api.ListEvents: Create filter")
	filter := hermes.Filter{
		Source:       req.FormValue("source"),
		ResourceType: req.FormValue("resource_type"),
		ResourceName: req.FormValue("resource_name"),
		UserName:     req.FormValue("user_name"),
		EventType:    req.FormValue("event_type"),
		Time:         req.FormValue("time"),
		Offset:       offset,
		Limit:        limit,
		Sort:         req.FormValue("sort"),
	}

	util.LogDebug("api.ListEvents: call hermes.GetEvents()")
	tenantId, err := getTenantId(req, res)
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
	protocol := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	// Do we need a NextURL?
	if int(filter.Offset+filter.Limit) < total {
		req.Form.Set("offset", strconv.FormatUint(filter.Offset+filter.Limit, 10))
		eventList.NextURL = fmt.Sprintf("%s://%s%s?%s", protocol, req.Host, req.URL.Path, req.Form.Encode())
	}
	// Do we need a PrevURL?
	if int(filter.Offset-filter.Limit) >= 0 {
		req.Form.Set("offset", strconv.FormatUint(filter.Offset-filter.Limit, 10))
		eventList.PrevURL = fmt.Sprintf("%s://%s%s?%s", protocol, req.Host, req.URL.Path, req.Form.Encode())
	}

	ReturnJSON(res, 200, eventList)
}

//GetEvent handles GET /v1/events/:event_id.
func (p *v1Provider) GetEventDetails(res http.ResponseWriter, req *http.Request) {
	token := p.CheckToken(req)
	if !token.Require(res, "event:show") {
		return
	}
	eventID := mux.Vars(req)["event_id"]
	tenantId, err := getTenantId(req, res)
	if err != nil {
		return
	}

	event, err := hermes.GetEvent(eventID, tenantId, p.keystone, p.storage)

	if ReturnError(res, err) {
		return
	}
	if event == nil {
		err := errors.New(fmt.Sprintf("Event %s could not be found in tenant %s", eventID, tenantId))
		http.Error(res, err.Error(), 404)
		return
	}
	ReturnJSON(res, 200, event)
}

func getTenantId(r *http.Request, w http.ResponseWriter) (string, error) {
	projectId := r.FormValue("project_id")
	domainId := r.FormValue("domain_id")
	var tenantId string
	if projectId != "" {
		tenantId = projectId
	}
	if domainId != "" {
		if projectId != "" {
			err := errors.New("domain_id and project_id cannot both be specified")
			http.Error(w, err.Error(), 400)
			return "", err
		} else {
			tenantId = domainId
		}
	}
	return tenantId, nil
}
