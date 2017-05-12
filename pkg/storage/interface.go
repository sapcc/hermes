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

import (
	"github.com/sapcc/hermes/pkg/data"
	"github.com/spf13/viper"
	"log"
)

// Driver is an interface that wraps the underlying event storage mechanism.
// Because it is an interface, the real implementation can be mocked away in unit tests.
type Interface interface {

	/********** requests to Keystone **********/
	GetEvents(filter data.Filter, tenant_id string) ([]*data.Event, int, error)
	GetEvent(eventId string, tenant_id string) (data.EventDetail, error)
}

func ConfiguredDriver() Interface {
	driverName := viper.GetString("hermes.storage_driver")
	switch driverName {
	case "elasticsearch":
		return ElasticSearch()
	case "mock":
		return Mock()
	default:
		log.Printf("Couldn't match a storage driver for configured value \"%s\"", driverName)
		return nil
	}
}
