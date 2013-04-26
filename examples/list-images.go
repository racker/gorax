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

	images, err := region.Images()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-36s   Name\n", "UUID")
	for _, i := range images {
		fmt.Printf("%36s - %s\n", string(i.Id[0:36]), i.Name)
	}
}
