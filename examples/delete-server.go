// vim: ts=8 sw=8 noet ai

package main

import (
	"flag"
	"github.com/racker/gorax/v2.0/cloud/servers"
	"github.com/racker/gorax/v2.0/identity"
	"log"
)

var pUserName = flag.String("u", "", "Rackspace API username")
var pPassword = flag.String("p", "", "Rackspace API password")
var pId = flag.String("i", "", "Server ID to delete")
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
		log.Fatal("You must specify the server ID to delete with the -i flag.")
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

	err = region.DeleteServerById(*pId)
	if err != nil {
		log.Fatal(err)
	}
}
