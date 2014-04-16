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
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/etcd/third_party/github.com/coreos/go-log/log"
	"github.com/racker/gorax"
	"github.com/racker/gorax/identity"
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

func (m *MonitoringClient) DeleteEntity(entityId string) (*Entity, error) {
	restReq := &gorax.RestRequest{
		Method:              "DELETE",
		Path:                "/entities/" + entityId,
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	_, err := m.client.PerformRequest(restReq)

	if err != nil {
		return nil, err
	}

	return nil, err
}

func (m *MonitoringClient) HostInfoEntity(entityId string, hostInfoType string) (interface{}, error) {
	var info interface{}

	switch hostInfoType {
	case "cpus":
		info = &CpuHostInfo{}
	case "memory":
		info = &MemoryHostInfo{}
	case "network_interfaces":
		info = &NetworkInterfaceHostInfo{}
	case "system":
		info = &SystemHostInfo{}
	case "disks":
		info = &DiskHostInfo{}
	case "filesystems":
		info = &FilesystemsHostInfo{}
	case "processes":
		info = &ProcessesHostInfo{}
	default:
		log.Error("Invalid Type")
		os.Exit(1)
	}

	path := fmt.Sprintf("/entities/%s/agent/host_info/%s", entityId, hostInfoType)
	restReq := &gorax.RestRequest{
		Method:              "GET",
		Path:                path,
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	resp, err := m.client.PerformRequest(restReq)
	if err != nil {
		return nil, err
	}
	resp.DeserializeBody(info)

	return info, err
}

func (m *MonitoringClient) AgentTargets(entityId string, agentType string) (interface{}, error) {
	info := &AgentTarget{}

	path := fmt.Sprintf("/entities/%s/agent/check_types/%s/targets", entityId, agentType)
	restReq := &gorax.RestRequest{
		Method:              "GET",
		Path:                path,
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	resp, err := m.client.PerformRequest(restReq)
	if err != nil {
		return nil, err
	}
	resp.DeserializeBody(info)

	return info, err
}

func (m *MonitoringClient) AgentTokenList() ([]AgentToken, error) {
	tokens := make([]AgentToken, 0)
	var nextMarker *string

	for true {
		restReq := &gorax.RestRequest{
			Method:              "GET",
			Path:                "/agent_tokens",
			ExpectedStatusCodes: []int{http.StatusOK},
		}

		if nextMarker != nil {
			restReq.Path += "?marker=" + *nextMarker
		}

		resp, err := m.client.PerformRequest(restReq)

		if err != nil {
			return nil, err
		}

		container := &PaginatedAgentTokenList{}
		err = resp.DeserializeBody(container)

		if err != nil {
			return nil, err
		}

		tokens = append(tokens, container.Values...)

		if container.Metadata.NextMarker == nil {
			break
		} else {
			nextMarker = container.Metadata.NextMarker
		}
	}

	return tokens, nil
}

func (m *MonitoringClient) AgentHostInfo(agentId string, agentType string) (interface{}, error) {
	var info interface{}

	switch agentType {
	case "cpus":
		info = &CpuHostInfo{}
	case "memory":
		info = &MemoryHostInfo{}
	case "network_interfaces":
		info = &NetworkInterfaceHostInfo{}
	case "system":
		info = &SystemHostInfo{}
	case "disks":
		info = &DiskHostInfo{}
	case "filesystems":
		info = &FilesystemsHostInfo{}
	case "processes":
		info = &ProcessesHostInfo{}
	default:
		log.Error("Invalid Type")
		os.Exit(1)
	}

	path := fmt.Sprintf("/agents/%s/host_info/%s", agentId, agentType)
	restReq := &gorax.RestRequest{
		Method:              "GET",
		Path:                path,
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	resp, err := m.client.PerformRequest(restReq)
	if err != nil {
		return nil, err
	}
	resp.DeserializeBody(info)

	return info, err
}

func (m *MonitoringClient) AgentConnectionsList(agentId string) (interface{}, error) {
	conns := make([]AgentConnection, 0)
	var nextMarker *string

	for true {
		restReq := &gorax.RestRequest{
			Method:              "GET",
			Path:                fmt.Sprintf("/agents/%s/connections", agentId),
			ExpectedStatusCodes: []int{http.StatusOK},
		}

		if nextMarker != nil {
			restReq.Path += "?marker=" + *nextMarker
		}

		resp, err := m.client.PerformRequest(restReq)

		if err != nil {
			return nil, err
		}

		container := &PaginatedAgentConnectionList{}
		err = resp.DeserializeBody(container)

		if err != nil {
			return nil, err
		}

		conns = append(conns, container.Values...)

		if container.Metadata.NextMarker == nil {
			break
		} else {
			nextMarker = container.Metadata.NextMarker
		}
	}

	return conns, nil
}

func (m *MonitoringClient) CheckTypeList() (interface{}, error) {
	types := make([]CheckType, 0)
	var nextMarker *string

	for true {
		restReq := &gorax.RestRequest{
			Method:              "GET",
			Path:                "/check_types",
			ExpectedStatusCodes: []int{http.StatusOK},
		}

		if nextMarker != nil {
			restReq.Path += "?marker=" + *nextMarker
		}

		resp, err := m.client.PerformRequest(restReq)

		if err != nil {
			return nil, err
		}

		container := &PaginatedCheckTypeList{}
		err = resp.DeserializeBody(container)

		if err != nil {
			return nil, err
		}

		types = append(types, container.Values...)

		if container.Metadata.NextMarker == nil {
			break
		} else {
			nextMarker = container.Metadata.NextMarker
		}
	}

	return types, nil
}

func (m *MonitoringClient) DeleteAgentToken(id string) error {
	path := fmt.Sprintf("/agent_tokens/%s", id)
	restReq := &gorax.RestRequest{
		Method:              "DELETE",
		Path:                path,
		ExpectedStatusCodes: []int{http.StatusNoContent},
	}

	_, err := m.client.PerformRequest(restReq)
	if err != nil {
		return err
	}

	return nil
}

func (m *MonitoringClient) ListMonitoringZones() (interface{}, error) {
	zones := make([]MonitoringZone, 0)
	var nextMarker *string

	for true {
		restReq := &gorax.RestRequest{
			Method:              "GET",
			Path:                "/monitoring_zones",
			ExpectedStatusCodes: []int{http.StatusOK},
		}

		if nextMarker != nil {
			restReq.Path += "?marker=" + *nextMarker
		}

		resp, err := m.client.PerformRequest(restReq)

		if err != nil {
			return nil, err
		}

		container := &PaginatedMonitoringZoneList{}
		err = resp.DeserializeBody(container)

		if err != nil {
			return nil, err
		}

		zones = append(zones, container.Values...)

		if container.Metadata.NextMarker == nil {
			break
		} else {
			nextMarker = container.Metadata.NextMarker
		}
	}

	return zones, nil

}

func (m *MonitoringClient) TracerouteMonitoringZone(mzId string, target string, resolver string) (interface{}, error) {
	postData := struct {
		Target         string `json:"target"`
		TargetResolver string `json:"target_resolver"`
	}{
		target,
		resolver,
	}

	path := fmt.Sprintf("/monitoring_zones/%s/traceroute", mzId)
	restReq := &gorax.RestRequest{
		Method: "POST",
		Path:   path,
		Body: &gorax.JSONRequestBody{
			Object: postData,
		},
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	route := &MonitoringZoneTraceroute{}

	resp, err := m.client.PerformRequest(restReq)
	if err != nil {
		return nil, err
	}
	resp.DeserializeBody(route)

	return route, err
}

func (m *MonitoringClient) ListMetrics(enId string, chId string) (interface{}, error) {
	metrics := make([]Metric, 0)
	var nextMarker *string

	path := fmt.Sprintf("/entities/%s/checks/%s/metrics", enId, chId)

	for true {
		restReq := &gorax.RestRequest{
			Method:              "GET",
			Path:                path,
			ExpectedStatusCodes: []int{http.StatusOK},
		}

		if nextMarker != nil {
			restReq.Path += "?marker=" + *nextMarker
		}

		resp, err := m.client.PerformRequest(restReq)

		if err != nil {
			return nil, err
		}

		container := &PaginatedMetricList{}
		err = resp.DeserializeBody(container)

		if err != nil {
			return nil, err
		}

		metrics = append(metrics, container.Values...)

		if container.Metadata.NextMarker == nil {
			break
		} else {
			nextMarker = container.Metadata.NextMarker
		}
	}

	return metrics, nil
}

func (m *MonitoringClient) ListLimits() (interface{}, error) {

	restReq := &gorax.RestRequest{
		Method:              "GET",
		Path:                "/limits",
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	limit := &Limit{}

	resp, err := m.client.PerformRequest(restReq)
	if err != nil {
		return nil, err
	}
	resp.DeserializeBody(limit)

	return limit, err
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
