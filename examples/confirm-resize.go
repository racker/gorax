// vim: ts=8 sw=8 noexpandtab ai

package main

import (
	"flag"
	"fmt"
	"github.com/racker/gorax/v2.0/cloud/servers"
	"github.com/racker/gorax/v2.0/identity"
	"log"
	"time"
)

var userName = flag.String("u", "", "Rackspace account username (required)")
var passWord = flag.String("p", "", "Rackspace account password (required)")
var region = flag.String("r", "DFW", "Rackspace region in which to create the server")
var serverId = flag.String("i", "", "ID of server to resize (required)")
var wait = flag.Bool("w", false, "Wait for server to become ready for confirmation first")

func waitForServer(region servers.Region, id string) error {
	ok := map[string]bool{
		"VERIFY_RESIZE": true,
		"ACTIVE":        true,
		"ERROR":         true,
	}
	for {
		s, err := region.ServerInfoById(id)
		if err != nil {
			return err
		}
		if ok[s.Status] {
			log.Printf("%s\n", s.Status)
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	panic("impossible")
}

func main() {
	flag.Parse()

	validations := map[string]string{
		"a username (-u flag)":  *userName,
		"a password (-p flag)":  *passWord,
		"a server ID (-i flag)": *serverId,
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

	if *wait {
		// TODO(sfalvo):
		//
		// I, or someone, will need to find out why this call is not
		// reliable.  On many occasions, I've observed that Rackspace
		// will set the server in ACTIVE status, despite it still being
		// in the rebuild operation per the mycloud.rackspace.com UI.
		// This will cause the call to ConfirmResizeServer() to respond
		// with a 409 error instead of a 204, as expected.  While
		// inconvenient, it doesn't appear to harm anything, so the
		// solution for now is to just re-run the command again.
		err = waitForServer(region, *serverId)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = region.ConfirmResizeServer(*serverId)
	if err != nil {
		log.Fatal(err)
	}
}
