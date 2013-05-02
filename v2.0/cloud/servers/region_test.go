// vim: ts=8 sw=8 noet ai

package servers

import (
	"github.com/racker/gorax/v2.0/identity"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const (
	SUCCESSFUL_LOGIN_RESPONSE = `{
	"access": {
		"serviceCatalog": [{
			"endpoints": [{
				"publicURL": "https://ord.servers.api.rackspacecloud.com/v2/12345",
				"region": "ORD",
				"tenantId": "12345",
				"versionId": "2",
				"versionInfo": "https://ord.servers.api.rackspacecloud.com/v2",
				"versionList": "https://ord.servers.api.rackspacecloud.com/"
			},{
				"publicURL": "https://dfw.servers.api.rackspacecloud.com/v2/12345",
				"region": "DFW",
				"tenantId": "12345",
				"versionId": "2",
				"versionInfo": "https://dfw.servers.api.rackspacecloud.com/v2",
				"versionList": "https://dfw.servers.api.rackspacecloud.com/"
			}],
			"name": "cloudServersOpenStack",
			"type": "compute"
		},{
			"endpoints": [{
				"publicURL": "https://ord.databases.api.rackspacecloud.com/v1.0/12345",
				"region": "ORD",
				"tenantId": "12345"
			}],
			"name": "cloudDatabases",
			"type": "rax:database"
		}],
		"token": {
			"expires": "2012-04-13T13:15:00.000-05:00",
			"id": "aaaaa-bbbbb-ccccc-dddd"
		},
		"user": {
			"RAX-AUTH:defaultRegion": "DFW",
			"id": "161418",
			"name": "demoauthor",
			"roles": [{
				"description": "User Admin Role.",
				"id": "3",
				"name": "identity:user-admin"
			}]
		}
	}
}
`
	TWO_IMAGES = `{
	"images": [
		{
			"OS-DCF:diskConfig": "AUTO", 
			"created": "2012-10-13T16:53:56Z", 
			"id": "a3a2c42f-575f-4381-9c6d-fcd3b7d07d28", 
			"links": [
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/v2/658405/images/a3a2c42f-575f-4381-9c6d-fcd3b7d07d28", 
					"rel": "self"
				}, 
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/658405/images/a3a2c42f-575f-4381-9c6d-fcd3b7d07d28", 
					"rel": "bookmark"
				}, 
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/658405/images/a3a2c42f-575f-4381-9c6d-fcd3b7d07d28", 
					"rel": "alternate", 
					"type": "application/vnd.openstack.image"
				}
			], 
			"metadata": {
				"arch": "x86-64", 
				"auto_disk_config": "True", 
				"com.rackspace__1__build_core": "1", 
				"com.rackspace__1__build_managed": "0", 
				"com.rackspace__1__build_rackconnect": "1", 
				"com.rackspace__1__options": "0", 
				"com.rackspace__1__visible_core": "1", 
				"com.rackspace__1__visible_managed": "0", 
				"com.rackspace__1__visible_rackconnect": "1", 
				"image_type": "base", 
				"org.openstack__1__architecture": "x64", 
				"org.openstack__1__os_distro": "org.centos", 
				"org.openstack__1__os_version": "5.8", 
				"os_distro": "centos", 
				"os_type": "linux", 
				"os_version": "5.8", 
				"rax_managed": "false", 
				"rax_options": "0"
			}, 
			"minDisk": 10, 
			"minRam": 256, 
			"name": "CentOS 5.8", 
			"progress": 100, 
			"status": "ACTIVE", 
			"updated": "2012-10-13T16:54:55Z"
		}, 
		{
			"OS-DCF:diskConfig": "AUTO", 
			"created": "2012-10-13T16:53:56Z", 
			"id": "a3a2c42f-575f-4381-9c6d-fcd3b7d07d17", 
			"links": [
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/v2/658405/images/a3a2c42f-575f-4381-9c6d-fcd3b7d07d17", 
					"rel": "self"
				}, 
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/658405/images/a3a2c42f-575f-4381-9c6d-fcd3b7d07d17", 
					"rel": "bookmark"
				}, 
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/658405/images/a3a2c42f-575f-4381-9c6d-fcd3b7d07d17", 
					"rel": "alternate", 
					"type": "application/vnd.openstack.image"
				}
			], 
			"metadata": {
				"arch": "x86-64", 
				"auto_disk_config": "True", 
				"com.rackspace__1__build_core": "1", 
				"com.rackspace__1__build_managed": "0", 
				"com.rackspace__1__build_rackconnect": "1", 
				"com.rackspace__1__options": "0", 
				"com.rackspace__1__visible_core": "1", 
				"com.rackspace__1__visible_managed": "0", 
				"com.rackspace__1__visible_rackconnect": "1", 
				"image_type": "base", 
				"org.openstack__1__architecture": "x64", 
				"org.openstack__1__os_distro": "org.centos", 
				"org.openstack__1__os_version": "6.0", 
				"os_distro": "centos", 
				"os_type": "linux", 
				"os_version": "6.0", 
				"rax_managed": "false", 
				"rax_options": "0"
			}, 
			"minDisk": 10, 
			"minRam": 256, 
			"name": "CentOS 6.0", 
			"progress": 100, 
			"status": "ACTIVE", 
			"updated": "2012-10-13T16:54:55Z"
		}
	]
}`

	TWO_FLAVORS = `{
	"flavors": [
		{
			"OS-FLV-DISABLED:disabled": false,
			"disk": 40,
			"id": "3",
			"links": [
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/v2/010101/flavors/3",
					"ref": "self"
				}, {
					"href": "https://dfw.servers.api.rackspacecloud.com/010101/flavors/3",
					"ref": "bookmark"
				}
			],
			"name": "1GB Standard Instance",
			"ram": 1024, 
			"rxtx_factor": 3.0, 
			"swap": 1024,
			"vcpus": 1
		}, {
			"OS-FLV-DISABLED:disabled": false,
			"disk": 80,
			"id": "4",
			"links": [
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/v2/010101/flavors/4",
					"ref": "self"
				}, {
					"href": "https://dfw.servers.api.rackspacecloud.com/010101/flavors/4",
					"ref": "bookmark"
				}
			],
			"name": "2GB Standard Instance",
			"ram": 2048, 
			"rxtx_factor": 3.0, 
			"swap": 1024,
			"vcpus": 1
		}
	]
}`
)

// testTransport is used to intercept traffic that would normally go out over a network connection.
// For the purposes of this package, we're concerned with the following things:
//
// The response string substitutes for the server response.  Setting this field, we can control
// what a test sees at any given time, allowing us to fake both error and successful conditions
// in full isolation of any provided network.
//
// The seenXAuthToken field records whether or not an X-Auth-Token has been provided by the client.
// Since we require an authenticated identity to access region-provided services,
// this header must always be present.
type testTransport struct {
	response       string
	seenXAuthToken bool
}

// The RoundTrip method implements the net/http.RoundTripper interface.
// It's here that we wrest control from the normal network stack and inject our own
// responses to the net/http package.
func (t *testTransport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	if req.Header.Get("X-Auth-Token") != "" {
		t.seenXAuthToken = true
	}

	headers := make(http.Header)
	body := ioutil.NopCloser(strings.NewReader(t.response))
	rsp = &http.Response{
		Status:           "200 OK",
		StatusCode:       200,
		Proto:            "HTTP/1.0",
		ProtoMajor:       1,
		ProtoMinor:       0,
		Header:           headers,
		Body:             body,
		ContentLength:    -1,
		TransferEncoding: nil,
		Close:            true,
		Trailer:          nil,
		Request:          req,
	}
	return
}

// withTestTransport abstracts common set-up code for creating a custom transport.
// The custom HTTP transport allows us to intercept normal HTTP interactions and
// fake out all network activity as we see fit.
func withTestTransport(r string, f func(c *http.Client, t *testTransport)) {
	transport := &testTransport{
		response: r,
	}
	client := &http.Client{
		Transport: transport,
	}
	f(client, transport)
}

// withAuthentication abstracts common set-up code for authenticating an identity.
func withAuthentication(c *http.Client, f func(e error, id identity.Identity)) {
	id := identity.NewIdentity("unused", "fields", "")
	id.UseClient(c)
	err := id.Authenticate()
	if err != nil {
		f(err, nil)
		return
	}
	f(nil, id)
}

// withRegion finds a known region and
// configures the transport for the (anticipated) next step in the testing process.
func withRegion(err error, id identity.Identity, cl *http.Client, t *testTransport, r string, f func(err error, _ Region)) {
	if err != nil {
		f(err, nil)
		return
	}
	region, err := RegionByName(id, "dfw")
	if err != nil {
		f(err, nil)
		return
	}
	region.UseClient(cl)
	t.response = r
	f(nil, region)
}

/****** Unit Tests ******/

func TestEndpointByName(t *testing.T) {
	withTestTransport(SUCCESSFUL_LOGIN_RESPONSE, func(client *http.Client, transport *testTransport) {
		withAuthentication(client, func(err error, id identity.Identity) {
			withRegion(err, id, client, transport, "", func(err error, region Region) {
				if err != nil {
					t.Error(err)
					return
				}
				api, err := region.EndpointByName("images")
				if err != nil {
					t.Error(err)
					return
				}
				if api != "https://dfw.servers.api.rackspacecloud.com/v2/12345/images" {
					t.Error("Expected DFW cloud server API for images; got", api)
					return
				}

				api, err = region.EndpointByName("flavors")
				if err != nil {
					t.Error(err)
					return
				}
				if api != "https://dfw.servers.api.rackspacecloud.com/v2/12345/flavors" {
					t.Error("Expected DFW cloud server API for flavors; got", api)
					return
				}

				api, err = region.EndpointByName("servers")
				if err != nil {
					t.Error(err)
					return
				}
				if api != "https://dfw.servers.api.rackspacecloud.com/v2/12345/servers" {
					t.Error("Expected DFW cloud server API for servers; got", api)
					return
				}
			})
		})
	})
}

func TestImages(t *testing.T) {
	withTestTransport(SUCCESSFUL_LOGIN_RESPONSE, func(client *http.Client, transport *testTransport) {
		withAuthentication(client, func(err error, id identity.Identity) {
			withRegion(err, id, client, transport, TWO_IMAGES, func(err error, region Region) {
				if err != nil {
					t.Error(err)
					return
				}
				imgs, err := region.Images()
				if err != nil {
					t.Error(err)
					return
				}
				if len(imgs) != 2 {
					t.Error("Expected 2 images; got", len(imgs))
					return
				}
				if !transport.seenXAuthToken {
					t.Error("Expected X-Auth-Token header to be sent")
					return
				}
			})
		})
	})
}

func TestFlavors(t *testing.T) {
	withTestTransport(SUCCESSFUL_LOGIN_RESPONSE, func(client *http.Client, transport *testTransport) {
		withAuthentication(client, func(err error, id identity.Identity) {
			withRegion(err, id, client, transport, TWO_FLAVORS, func(err error, region Region) {
				if err != nil {
					t.Error(err)
					return
				}
				flavors, err := region.Flavors()
				if err != nil {
					t.Error(err)
					return
				}
				if len(flavors) != 2 {
					t.Error("Expected 2 flavors; got", len(flavors))
					return
				}
			})
		})
	})
}
