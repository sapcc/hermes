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
	"errors"
	policy "github.com/databus23/goslo.policy"
	"github.com/gophercloud/gophercloud"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

//Token represents a user's token, as passed through the X-Auth-Token header of
//a request.
type Token struct {
	enforcer *policy.Enforcer
	context  policy.Context
	err      error
}

//CheckToken checks the validity of the request's X-Auth-Token in keystone, and
//returns a Token instance for checking authorization. Any errors that occur
//during this function are deferred until Require() is called.
func (p *v1Provider) CheckToken(r *http.Request) *Token {
	str := r.Header.Get("X-Auth-Token")
	if str == "" {
		return &Token{err: errors.New("X-Auth-Token header missing")}
	}

	t := &Token{enforcer: viper.Get("hermes.PolicyEnforcer").(*policy.Enforcer)}
	t.context, t.err = p.keystone.ValidateToken(str)
	switch t.err.(type) {
	case gophercloud.ErrDefault404:
		t.err = errors.New("X-Auth-Token is invalid or expired")
	}
	t.context.Request = mux.Vars(r)
	if r.FormValue("domain_id") == "" {
		t.context.Request["domain_id"] = t.context.Auth["domain_id"]
	} else {
		t.context.Request["domain_id"] = r.FormValue("domain_id")
	}
	if r.FormValue("project_id") == "" {
		t.context.Request["project_id"] = t.context.Auth["project_id"]
	} else {
		t.context.Request["project_id"] = r.FormValue("project_id")
	}
	return t
}

//Require checks if the given token has the given permission according to the
//policy.json that is in effect. If not, an error response is written and false
//is returned.
func (t *Token) Require(w http.ResponseWriter, rule string) bool {
	if t.err != nil {
		http.Error(w, t.err.Error(), 401)
		return false
	}

	if os.Getenv("DEBUG") == "1" {
		t.context.Logger = log.Printf //or any other function with the same signature
	}
	if !t.enforcer.Enforce(rule, t.context) {
		http.Error(w, "Unauthorized", 403)
		return false
	}
	return true
}
