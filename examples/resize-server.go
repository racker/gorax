// vim: ts=8 sw=8 noexpandtab ai

package main

import (
	"flag"
	"fmt"
	"github.com/racker/gorax/v2.0/cloud/servers"
	"github.com/racker/gorax/v2.0/identity"
	"log"
	"strings"
)

var userName = flag.String("u", "", "Rackspace account username (required)")
var passWord = flag.String("p", "", "Rackspace account password (required)")
var region = flag.String("r", "DFW", "Rackspace region in which to create the server")
var serverId = flag.String("i", "", "ID of server to resize (required)")
var flavorId = flag.String("f", "", "The flavor for the rebuilt server (required)")
var diskConfig = flag.String("d", "", "The new disk configuration value (optional); choices are '' for leaving the same, AUTO, or MANUAL.")
var serverName = flag.String("n", "", "The new server name (required)")

func main() {
	flag.Parse()

	validations := map[string]string{
		"a username (-u flag)":    *userName,
		"a password (-p flag)":    *passWord,
		"a flavor ID (-f flag)":   *flavorId,
		"a server name (-n flag)": *serverName,
		"a server ID (-i flag)":   *serverId,
	}
	for flag, value := range validations {
		if value == "" {
			log.Fatal(fmt.Sprintf("You must provide %s", flag))
		}
	}

	config := strings.ToUpper(*diskConfig)
	if (config != "") && (config != "AUTO") && (config != "MANUAL") {
		log.Printf("The disk configuration must be one of the following values:\n")
		log.Printf("  (unspecifed) - Leave it the same.")
		log.Printf("  AUTO         - The server is built with a single partition the size of the target flavor disk.")
		log.Fatal("  MANUAL       - The server is built using whatever partition scheme and file system is in the source image.")
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

	err = region.ResizeServer(*serverId, *serverName, *flavorId, config)
	if err != nil {
		log.Fatal(err)
	}
}
