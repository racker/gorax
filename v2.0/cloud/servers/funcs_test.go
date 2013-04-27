// vim: ts=8 sw=8 noexpandtab ai

package servers

import (
	"github.com/racker/gorax/v2.0/identity"
	"testing"
)

/****** Fake Identities ******/

// The myIdCard structure substitues for an authenticated identity that
// only offers legacy cloudServers, but no OpenStack servers.
type myIdCard struct{}

func (i *myIdCard) Password() string {
	return ""
}

func (i *myIdCard) SetCredentials(userName, key, reg string) {
}

func (i *myIdCard) Username() string {
	return "my-username"
}

func (i *myIdCard) Region() string {
	return "my-region"
}

func (i *myIdCard) Token() (string, error) {
	return "my-token", nil
}

func (i *myIdCard) Expires() (string, error) {
	return "2020-01-01T12:00:00", nil
}

func (i *myIdCard) TenantId() (string, error) {
	return "123456", nil
}

func (i *myIdCard) TenantName() (string, error) {
	return "tenant-name", nil
}

func (i *myIdCard) AuthEndpoint() (ep string) {
	return "http://example.com/auth-endpoint"
}

func (i *myIdCard) IsAuthenticated() bool {
	return true
}

func (i *myIdCard) ServiceCatalog() (sc []identity.CatalogEntry, err error) {
	sc = []identity.CatalogEntry{
		identity.CatalogEntry{
			Name: "cloudServers",
			Type: "compute",
			Endpoints: []identity.EntryEndpoint{
				identity.EntryEndpoint{
					TenantId:    "123456",
					PublicURL:   "http://servers.api.rackspacecloud.com/v1.0/123456",
					VersionId:   "1.0",
					VersionInfo: "https://servers.api.rackspacecloud.com/v1.0/",
					VersionList: "https://servers.api/rackspacecloud.com/",
				},
			},
		},
	}
	return
}

func (i *myIdCard) Roles() ([]identity.Role, error) {
	return nil, nil
}

func (i *myIdCard) Authenticate() error {
	return nil
}

// Constructor for myIdCard.
func fakeId() identity.Identity {
	return &myIdCard{}
}

// The myIdCard2 structure provides the same benefits as myIdCard,
// except that it also provides OpenStack cloud servers.
type myIdCard2 struct {
	myIdCard
}

func (i *myIdCard2) ServiceCatalog() (sc []identity.CatalogEntry, err error) {
	sc = []identity.CatalogEntry{
		identity.CatalogEntry{
			Name: "cloudServers",
			Type: "compute",
			Endpoints: []identity.EntryEndpoint{
				identity.EntryEndpoint{
					TenantId:    "123456",
					PublicURL:   "http://servers.api.rackspacecloud.com/v1.0/123456",
					VersionId:   "1.0",
					VersionInfo: "https://servers.api.rackspacecloud.com/v1.0/",
					VersionList: "https://servers.api/rackspacecloud.com/",
				},
			},
		},
		identity.CatalogEntry{
			Name: "cloudServersOpenStack",
			Type: "compute",
			Endpoints: []identity.EntryEndpoint{
				identity.EntryEndpoint{
					PublicURL:   "https://dfw.servers.api.rackspacecloud.com/v2/775360",
					Region:      "DFW",
					TenantId:    "775360",
					VersionId:   "2",
					VersionInfo: "https://dfw.servers.api.rackspacecloud.com/v2",
					VersionList: "https://dfw.servers.api.rackspacecloud.com/",
				},
				identity.EntryEndpoint{
					PublicURL:   "https://ord.servers.api.rackspacecloud.com/v2/775360",
					Region:      "ORD",
					TenantId:    "775360",
					VersionId:   "2",
					VersionInfo: "https://ord.servers.api.rackspacecloud.com/v2",
					VersionList: "https://ord.servers.api.rackspacecloud.com/",
				},
			},
		},
	}
	return
}

// Constructor for myIdCard2.
func fakeId2() identity.Identity {
	return &myIdCard2{}
}

/****** Unit Tests ******/

func TestInRegion(t *testing.T) {
	id := fakeId()

	_, err := RegionByName(id, "ussr")
	if err == nil {
		t.Error("InRegion: unsupported location should yield error")
		return
	}

	_, err = RegionByName(id, "ord")
	if err == nil {
		t.Error("InRegion: unsupported location should yield error")
		return
	}

	id2 := fakeId2()

	_, err = RegionByName(id2, "ussr")
	if err == nil {
		t.Error("InRegion: unsupported location should yield error")
		return
	}

	_, err = RegionByName(id2, "ord")
	if err != nil {
		t.Error("InRegion: supported region shouldn't yield error; got", err)
		return
	}

	_, err = RegionByName(id2, "DFW")
	if err != nil {
		t.Error("InRegion: Case should not be sensitive; got error", err)
		return
	}
}
