// vim: ts=8 sw=8 noexpandtab

package main

import (
	"fmt"
	"github.com/racker/gorax/v2.0/cloud/servers"
	"github.com/racker/gorax/v2.0/identity"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: I need both username and API key on CLI, in that order.")
	}
	username := os.Args[1]
	password := os.Args[2]

	id := identity.NewIdentity(username, password, "")
	err := id.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	region, err := servers.RegionByName(id, "dfw")
	if err != nil {
		log.Fatal(err)
	}

	servers, err := region.Servers()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID,Name\n")
	for _, i := range servers {
		fmt.Printf("%s,\"%s\"\n", i.Id, i.Name)
	}
}
