// vim: ts=8 sw=8 noet ai

package main

import (
	"log"
	"fmt"
	"flag"
	"github.com/racker/gorax/v2.0/identity"
	"github.com/racker/gorax/v2.0/cloud/servers"
)


var pUserName = flag.String("u", "", "Rackspace API username")
var pPassword = flag.String("p", "", "Rackspace API password")
var pId = flag.String("i", "", "Server ID to get info on")
var pRegion = flag.String("r", "DFW", "Region where the server is hosted")


func main() {
	var err error

	flag.Parse()

	if *pUserName == "" {
		log.Fatal("You must specify a username with the -u flag.")
	}
	if *pPassword == "" {
		log.Fatal("You must specify an API password with the -p flag.")
	}
	if *pId == "" {
		log.Fatal("You must specify the server ID to get information on with the -i flag.")
	}

	id := identity.NewIdentity(*pUserName, *pPassword, *pRegion)
	err = id.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	region, err := servers.RegionByName(id, *pRegion)
	if err != nil {
		log.Fatal(err)
	}

	s, err := region.ServerInfoById(*pId)
	if err != nil {
		log.Fatal(err)
	}

	fields := map[string]string{
		"Access IPv4: %s": s.AccessIPv4,
		"Access IPv6: %s": s.AccessIPv6,
		"Created: %s": s.Created,
		"Flavor: %s": s.Flavor.Id,
		"Host ID: %s": s.HostId,
		"ID: %s": s.Id,
		"Image: %s": s.Image.Id,
		"Name: %s": s.Name,
		"Progress: %s": fmt.Sprintf("%d", s.Progress),
		"Status: %s": s.Status,
		"Tenant ID: %s": s.TenantId,
		"Updated: %s": s.Updated,
		"User ID: %s": s.UserId,
	}

	fmt.Printf("Server info:\n")
	for k, v := range fields {
		fmt.Printf(k, v)
		fmt.Println()
	}
}

