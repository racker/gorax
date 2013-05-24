// vim: ts=8 sw=8 noexpandtab ai

package main

import (
	"flag"
	"fmt"
	"github.com/racker/gorax/v2.0/cloud/servers"
	"github.com/racker/gorax/v2.0/identity"
	"log"
)

var userName = flag.String("u", "", "Rackspace account username")
var passWord = flag.String("p", "", "Rackspace account password")
var region = flag.String("r", "DFW", "Rackspace region in which to create the server")
var adminPassword = flag.String("w", "", "The administrator password of the server")
var serverId = flag.String("i", "", "The ID of the server to rebuild")
var imageId = flag.String("I", "", "Image ID to apply to the server")
var flavorId = flag.String("f", "", "The flavor for the rebuilt server")
var serverName = flag.String("n", "", "The new name of the server")

func main() {
	flag.Parse()

	validations := map[string]string{
		"a username (-u flag)":                *userName,
		"a password (-p flag)":                *passWord,
		"a server name (-n flag)":             *serverName,
		"an image ID (-I flag)":               *imageId,
		"a flavor ID (-f flag)":               *flavorId,
		"an administrator password (-w flag)": *adminPassword,
	}
	for flag, value := range validations {
		if value == "" {
			log.Fatal(fmt.Sprintf("You must provide %s", flag))
		}
	}

	id := identity.NewIdentity(*userName, *passWord, *region)
	err := id.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	region, err := servers.RegionByName(id, *region)
	if err != nil {
		log.Fatal(err)
	}

	_, err := region.RebuildServer(*serverId, servers.NewServer{
		Name:      *serverName,
		ImageRef:  *imageId,
		FlavorRef: *flavorId,
		AdminPass: *adminPassword,
	})
	if err != nil {
		log.Fatal(err)
	}
}
