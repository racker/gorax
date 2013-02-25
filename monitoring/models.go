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

type PaginationMetadata struct {
	Count      int     `json:"count"`
	Limit      int     `json:"limit"`
	Marker     *string `json:"marker"`
	NextMarker *string `json:"next_marker"`
	NextHref   *string `json:"next_href"`
}

type Entity struct {
	Id          string            `json:"id"`
	Label       *string           `json:"label"`
	Metadata    map[string]string `json:"metadata"`
	Managed     bool              `json:"managed"`
	Uri         *string           `json:"uri"`
	AgentId     *string           `json:"agent_id"`
	IPAddresses map[string]string `json:"ip_addresses"`
}

type PaginatedEntityList struct {
	Values   []Entity
	Metadata PaginationMetadata
}

type Check struct {
	Id                  string                 `json:"id"`
	Label               *string                `json:"label"`
	Type                string                 `json:"type"`
	Details             map[string]interface{} `json:"details"`
	MonitoringZonesPoll []string               `json:"monitoring_zones_poll"`
	Timeout             int                    `json:"timeout"`
	Period              int                    `json:"period"`
	TargetAlias         *string                `json:"target_alias"`
	TargetHostname      *string                `json:"target_hostname"`
	TargetResolver      *string                `json:"target_resolver"`
	Disabled            bool                   `json:"disabled"`
	Metadata            map[string]interface{} `json:"metadata"`
}

type PaginatedCheckList struct {
	Values   []Check
	Metadata PaginationMetadata
}
