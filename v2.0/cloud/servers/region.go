// vim: ts=8 sw=8 noet ai

package servers

import (
	"fmt"
	"github.com/racker/gorax/v2.0/identity"
	"github.com/racker/perigee"
	"net/http"
)

// A raxRegion represents a Rackspace-hosted region.
type raxRegion struct {
	id            identity.Identity
	entryEndpoint identity.EntryEndpoint
	httpClient    *http.Client
	token         string
}

// Flavors method provides a complete list of machine configurations (called flavors) available at the region.
func (r *raxRegion) Flavors() ([]Flavor, error) {
	var fc *FlavorsContainer
	var fs []Flavor

	url, _ := r.EndpointByName("flavors")
	err := perigee.Get(url, perigee.Options{
		CustomClient: r.httpClient,
		Results:      &fc,
		MoreHeaders: map[string]string{
			"X-Auth-Token": r.token,
		},
	})
	if err == nil {
		fs = fc.Flavors
	}
	return fs, err
}

// Images method provides a complete list of images hosted at the region.
func (r *raxRegion) Images() ([]Image, error) {
	var ic *ImagesContainer
	var is []Image

	url, _ := r.EndpointByName("images")
	err := perigee.Get(url, perigee.Options{
		CustomClient: r.httpClient,
		Results:      &ic,
		MoreHeaders: map[string]string{
			"X-Auth-Token": r.token,
		},
	})
	if err == nil {
		is = ic.Images
	}
	return is, err
}

// Servers method provides a complete list of servers hosted by the user
// at a given region.
func (r *raxRegion) Servers() ([]Server, error) {
	var sc *ServersContainer
	var ss []Server

	url, _ := r.EndpointByName("servers")
	err := perigee.Get(url, perigee.Options{
		CustomClient: r.httpClient,
		Results:      &sc,
		MoreHeaders: map[string]string{
			"X-Auth-Token": r.token,
		},
	})
	if err == nil {
		ss = sc.Servers
	}
	return ss, err
}

// CreateServer requests a new server to be created by the cloud provider.
// The user must pass in a pointer to an initialized NewServerContainer structure.
// Please refer to NewServerContainer for more details.
//
// If the NewServer structure lacks a field assignment for AdminPass,
// a password will be automatically generated by OpenStack / Rackspace and
// returned back through the AdminPass field.  Take care; this will be
// the only time this happens; no other means exists in the public API
// to acquire a password for a pre-existing server.
func (r *raxRegion) CreateServer(ns NewServer) (*NewServer, error) {
	var ns *NewServer

	nsc := &NewServerContainer{
		Server: ns,
	}
	ep, err := r.EndpointByName("servers")
	if err != nil {
		return nil, err
	}
	nsr := new(NewServerContainer)
	err = perigee.Post(ep, perigee.Options{
		ReqBody: nsc,
		Results: nsr,
		MoreHeaders: map[string]string{
			"X-Auth-Token": r.token,
		},
		OkCodes: []int{202},
	})
	if err == nil {
		ns = &nsr.Server
	}
	return ns, err
}

// ServerInfoById provides the complete server information record
// given you know its unique ID.
func (r *raxRegion) ServerInfoById(id string) (*Server, error) {
	var serverEnvelop *ServerInfoContainer
	var s *Server

	baseUrl, err := r.EndpointByName("servers")
	serverUrl := fmt.Sprintf("%s/%s", baseUrl, id)
	err = perigee.Get(serverUrl, perigee.Options{
		Results: &serverEnvelop,
		MoreHeaders: map[string]string {
			"X-Auth-Token": r.token,
		},
	})
	if err == nil {
		s = &serverEnvelop.Server
	}
	return s, err
}

// EndpointByName computes a resource URL, assuming a valid name.
// An error is returned if an invalid or unsupported endpoint name is given.
//
// It is an error for application software to invoke this method.
// This method exists and is publicly available only to support testing.
func (r *raxRegion) EndpointByName(name string) (string, error) {
	var supportedEndpoint map[string]bool = map[string]bool{
		"images":  true,
		"flavors": true,
		"servers": true,
	}

	if supportedEndpoint[name] {
		api := fmt.Sprintf("%s/%s", r.entryEndpoint.PublicURL, name)
		return api, nil
	}
	return "", fmt.Errorf("Unsupported endpoint")
}

// UseClient configures the region client to use a specific net/http client.
// This allows you to configure a custom HTTP transport for specialized requirements.
// You normally wouldn't need to set this, as the net/http package makes reasonable
// choices on its own.  Customized transports are useful, however, if extra logging
// is required, or if you're using unit tests to isolate and verify correct behavior.
func (r *raxRegion) UseClient(cl *http.Client) {
	r.httpClient = cl
}

// makeRegionalClient creates a structure that implements the Region interface.
func makeRegionalClient(id identity.Identity, e identity.EntryEndpoint) (Region, error) {
	t, err := id.Token()
	if err != nil {
		return nil, err
	}
	return &raxRegion{
		id:            id,
		entryEndpoint: e,
		token:         t,
		httpClient:    &http.Client{},
	}, nil
}
