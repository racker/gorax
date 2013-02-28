/*
Copyright 2013 Rackspace

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS-IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package monitoring

import (
	"github.com/racker/gorax"
	"github.com/racker/gorax/identity"
	"net/http"
)

// A MonitoringClient object exists for each outstanding connection to the Rackspace Cloud Monitoring APIs.
type MonitoringClient struct {
	client *gorax.RestClient
}

// SetDebug() configures whether or not the monitoring client works in debug-mode (true) or not (false).
func (m *MonitoringClient) SetDebug(debug bool) {
	m.client.SetDebug(debug)
}

func (m *MonitoringClient) ListEntities() ([]Entity, error) {
	entities := make([]Entity, 0)
	var nextMarker *string

	for true {
		restReq := &gorax.RestRequest{
			Method:              "GET",
			Path:                "/entities",
			ExpectedStatusCodes: []int{http.StatusOK},
		}

		if nextMarker != nil {
			restReq.Path += "?marker=" + *nextMarker
		}

		resp, err := m.client.PerformRequest(restReq)

		if err != nil {
			return nil, err
		}

		container := &PaginatedEntityList{}
		err = resp.DeserializeBody(container)

		if err != nil {
			return nil, err
		}

		entities = append(entities, container.Values...)

		if container.Metadata.NextMarker == nil {
			break
		} else {
			nextMarker = container.Metadata.NextMarker
		}
	}

	return entities, nil
}

func (m *MonitoringClient) GetEntity(entityId string) (*Entity, error) {
	restReq := &gorax.RestRequest{
		Method:              "GET",
		Path:                "/entities/" + entityId,
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	resp, err := m.client.PerformRequest(restReq)

	if err != nil {
		return nil, err
	}

	entity := &Entity{}
	err = resp.DeserializeBody(entity)
	return entity, err
}

// ListChecks() retrieves a list of Check objects configured for a given entity.
// This function abstracts pagination of the results for you.
// If successful, the error result will always be nil; otherwise, the Check slice will be nil.
func (m *MonitoringClient) ListChecks(entityId string) ([]Check, error) {
	checks := make([]Check, 0)
	var nextMarker *string

	for {
		restReq := &gorax.RestRequest{
			Method:              "GET",
			Path:                "/entities/" + entityId + "/checks",
			ExpectedStatusCodes: []int{http.StatusOK},
		}

		if nextMarker != nil {
			restReq.Path += "?marker=" + *nextMarker
		}

		resp, err := m.client.PerformRequest(restReq)

		if err != nil {
			return nil, err
		}

		container := &PaginatedCheckList{}
		err = resp.DeserializeBody(container)

		if err != nil {
			return nil, err
		}

		checks = append(checks, container.Values...)

		if container.Metadata.NextMarker == nil {
			break
		} else {
			nextMarker = container.Metadata.NextMarker
		}
	}

	return checks, nil
}

// MakePasswordMonitoringClient creates an object representing the monitoring client, with username/password authentication.
func MakePasswordMonitoringClient(url string, authurl string, username string, password string) *MonitoringClient {
	m := &MonitoringClient{
		client: gorax.MakeRestClient(url),
	}
	m.client.RequestMiddlewares = []gorax.RequestMiddleware{
		identity.MakeKeystonePasswordMiddleware(authurl, username, password),
	}
	return m
}

func MakeAPIKeyMonitoringClient(url string, authurl string, username string, apiKey string) *MonitoringClient {
	m := &MonitoringClient{
		client: gorax.MakeRestClient(url),
	}
	m.client.RequestMiddlewares = []gorax.RequestMiddleware{
		identity.MakeKeystoneAPIKeyMiddleware(authurl, username, apiKey),
	}
	return m
}
