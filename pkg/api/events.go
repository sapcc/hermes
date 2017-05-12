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
	"github.com/sapcc/hermes/pkg/data"
	"github.com/sapcc/hermes/pkg/hermes"
	"github.com/sapcc/hermes/pkg/storage"
	"strconv"
)

//ListEvents handles GET /v1/events.
func (p *v1Provider) ListEvents(w http.ResponseWriter, r *http.Request) {
	token := p.CheckToken(r)
	if !p.CheckToken(r).Require(w, "event:list") {
		return
	}

	source := r.FormValue("source")
	resource_type := r.FormValue("resource_type")
	resource_name := r.FormValue("resource_name")
	user_name := r.FormValue("user_name")
	event_type := r.FormValue("event_type")
	time := r.FormValue("time")
	offset, _ := strconv.ParseUint(r.FormValue("offset"), 10, 32)
	limit, _ := strconv.ParseUint(r.FormValue("limit"), 10, 8)
	sort := r.FormValue("sort")

	filter := data.Filter{
		Source:       source,
		ResourceType: resource_type,
		ResourceName: resource_name,
		UserName:     user_name,
		EventType:    event_type,
		Time:         time,
		Offset:       offset,
		Limit:        limit,
		Sort:         sort,
	}

	events, total, err := hermes.GetEvents(&filter, &token.context, storage.ConfiguredDriver())

	if ReturnError(w, err) {
		return
	}

	eventList := data.EventList{Events: events, Total: total}

	protocol := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	if int(filter.Offset+filter.Limit) < total {
		r.Form.Set("offset", strconv.FormatUint(filter.Offset+filter.Limit, 10))
		eventList.NextURL = fmt.Sprintf("%s://%s%s?%s", protocol, r.Host, r.URL.Path, r.Form.Encode())
	}
	if int(filter.Offset-filter.Limit) >= 0 {
		r.Form.Set("offset", strconv.FormatUint(filter.Offset-filter.Limit, 10))
		eventList.PrevURL = fmt.Sprintf("%s://%s%s?%s", protocol, r.Host, r.URL.Path, r.Form.Encode())
	}

	ReturnJSON(w, 200, eventList)
}

//GetEvent handles GET /v1/events/:event_id.
func (p *v1Provider) GetEventDetails(w http.ResponseWriter, r *http.Request) {
	token := p.CheckToken(r)
	if !p.CheckToken(r).Require(w, "event:show") {
		return
	}

	eventID := mux.Vars(r)["event_id"]

	event, err := hermes.GetEvent(eventID, &token.context, storage.ConfiguredDriver())
	if ReturnError(w, err) {
		return
	}

	ReturnJSON(w, 200, event)
}
