// vim: ts=8 sw=8 noet ai

package servers

import (
	"encoding/json"
	"fmt"
	"github.com/racker/gorax/v2.0/identity"
	"io/ioutil"
	"net/http"
)

// A raxRegion represents a Rackspace-hosted region.
type raxRegion struct {
	id            identity.Identity
	entryEndpoint identity.EntryEndpoint
	httpClient    *http.Client
	token         string
}

// ImagesContainer is used for JSON (un)marshalling.
// It provides the top-most container for image records.
type ImagesContainer struct {
	Images []Image `json:"images"`
}

// ImageLink is used for JSON (un)marshalling.
// It provides RESTful links to an image resource.
type ImageLink struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

// Image is used for JSON (un)marshalling.
// It provides a description of an OS image.
//
// The Id field contains the image's unique identifier.
// For example, this identifier will be useful for specifying which operating system to install on a new server instance.
//
// The MinDisk and MinRam fields specify the minimum resources a server must provide to be able to install the image.
//
// The Name field provides a human-readable moniker for the OS image.
//
// The Progress and Status fields indicate image-creation status.
// Any usable image will have 100% progress.
//
// The Updated field indicates the last time this image was changed.
type Image struct {
	OS_DCF_diskConfig string      `json:"OS-DCF:diskConfig"`
	Created           string      `json:"created"`
	Id                string      `json:"id"`
	Links             []ImageLink `json:"links"`
	MinDisk           int         `json:"minDisk"`
	MinRam            int         `json:"minRam"`
	Name              string      `json:"name"`
	Progress          int         `json:"progress"`
	Status            string      `json:"status"`
	Updated           string      `json:"updated"`
}

// Images method provides a complete list of images hosted at the region.
func (r *raxRegion) Images() ([]Image, error) {
	var images ImagesContainer

	apiUrl, _ := r.EndpointByName("images")
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Auth-Token", r.token)

	rsp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("200 OK expected; got %s", rsp.Status)
	}
	defer rsp.Body.Close()
	jsonContainer, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonContainer, &images)
	return images.Images, err
}

// EndpointByName computes a resource URL, assuming a valid name.
// An error is returned if an invalid or unsupported endpoint name is given.
//
// It is an error for application software to invoke this method.
// This method exists and is publicly available only to support testing.
func (r *raxRegion) EndpointByName(name string) (string, error) {
	var supportedEndpoint map[string]bool = map[string]bool{
		"images": true,
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
