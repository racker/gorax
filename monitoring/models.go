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

// A Check represents a check that the user configures on one of his or her entities.
// It contains the following fields:
//
// The Id field uniquely identifies the check amongst all others belonging to the user.
//
// The Label field identifies the check to a human operator.
//
// The Type field indicates what kind of check it is.
//
// The Details field provides a mapping of detail to detail-specific information.
//
// The MonitoringZonesPoll field provides a list of what to poll for this check.
//
// The Timeout field indicates how many seconds to wait for a response before the check fails.
//
// The Period field tells how frequently to perform the check, in seconds.
//
// The TargetAlias field does something; I just don't quite know what it is.
//
// The TargetHostname field identifies the host name of that which is being checked.
//
// TargetResolver field identifies the domain name resolver scoping the target hostname.
//
// The Disabled field is true if the check is disabled for the entity; false otherwise.
//
// The Metadata field provides a generic key/value store of miscellaneous bits of information relevant to this check.
// However, its implementation isn't very efficient at all.  This field is not intended for use as a general purpose
// key/value store.
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

// A PaginatedCheckList contains a finite subset of the complete set of checks a cloud monitoring user has configured.
// The Values field contains the array slice representing the set of Check objects.
type PaginatedCheckList struct {
	Values   []Check
	Metadata PaginationMetadata
}
