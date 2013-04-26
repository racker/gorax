// vim: ts=8 sw=8 noexpandtab ai

package servers

import (
	"fmt"
	"github.com/racker/gorax/v2.0/identity"
	"strings"
)

// RegionByName grants access to "region" in which a server may be created.
// Often, regions correspond closely with the physical data centers where servers are actually located.
//
// You need an authenticated identity to gain access to a region.
//
// Traditionally, regions are named for the closest corresponding airport code.
// For example, a data center located in Dallas, Texas will be known as "DFW".
// However, this needn't always be the case;
// one data center in the United Kingdom uses the name "LON", not "LHR".
//
// Region names are case insensitive for convenience; however,
// they're traditionally written in all uppercase letters.
func RegionByName(id identity.Identity, region string) (r Region, err error) {
	sc, err := id.ServiceCatalog()
	if err != nil {
		return nil, err
	}

	for _, entry := range sc {
		if entry.Type == "compute" {
			for _, endpoint := range entry.Endpoints {
				if strings.ToUpper(endpoint.Region) == strings.ToUpper(region) {
					return makeRegionalClient(id, endpoint)
				}
			}
		}
	}
	return nil, fmt.Errorf("Unsupported region or V1.0 services only")
}
