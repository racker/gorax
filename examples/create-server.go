// vim: ts=8 sw=8 noexpandtab

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/racker/gorax/v2.0/cloud/servers"
	"github.com/racker/gorax/v2.0/identity"
	"log"
)

var userName = flag.String("u", "", "Rackspace account username")
var passWord = flag.String("p", "", "Rackspace account password")
var region = flag.String("r", "DFW", "Rackspace region in which to create the server")
var serverName = flag.String("n", "", "Server name")
var imageRef = flag.String("i", "", "Image ID to deploy onto the server")
var flavorRef = flag.String("f", "", "Flavor of server to deploy image upon")
var adminPass = flag.String("a", "", "Administrator password (auto-assigned if none)")

func main() {
	flag.Parse()

	validations := map[string]string{
		"a username (-u flag)":         *userName,
		"a password (-p flag)":         *passWord,
		"a server name (-n flag)":      *serverName,
		"an image reference (-i flag)": *imageRef,
		"a server flavor (-f flag)":    *flavorRef,
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

	nsr, err := region.CreateServer(servers.NewServer{
		Name:      *serverName,
		ImageRef:  *imageRef,
		FlavorRef: *flavorRef,
		AdminPass: *adminPass,
	})
	if err != nil {
		log.Fatal(err)
	}

	nscJson, err := json.MarshalIndent(nsr, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(nscJson))

	servers, err := region.Servers()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID,Name\n")
	for _, i := range servers {
		fmt.Printf("%s,\"%s\"\n", i.Id, i.Name)
	}
}
