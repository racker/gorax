// vim: ts=8 sw=8 noexpandtab

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
var serverId = flag.String("i", "", "ID of server you want to reboot")
var hard = flag.Bool("hard", false, "Specify if you desire a hard reboot (true) versus a soft reboot (false; default)")

func main() {
	flag.Parse()

	validations := map[string]string{
		"a username (-u flag)":         *userName,
		"a password (-p flag)":         *passWord,
		"a server ID (-i flag)":	*serverId,
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

	err = region.RebootServer(*serverId, *hard)
	if err != nil {
		log.Fatal(err)
	}
}
