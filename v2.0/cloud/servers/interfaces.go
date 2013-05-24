// vim: ts=8 sw=8 noexpandtab ai

package servers

import (
	"net/http"
)

// A Region represents a geographical area with cloud computing resources.
type Region interface {
	Images() ([]Image, error)
	Flavors() ([]Flavor, error)
	Servers() ([]Server, error)
	CreateServer(NewServer) (*NewServer, error)
	ServerInfoById(string) (*Server, error)
	RebootServer(string, bool) error
	DeleteServerById(string) error
	SetAdminPassword(string, string) error
	RebuildServer(string, NewServer) (*Server, error)
	UseClient(*http.Client)
	EndpointByName(string) (string, error)
}
