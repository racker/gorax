// vim: ts=8 sw=8 noet ai

package identity

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const (
	USERNAME    = "joe_user"
	PASSWORD    = "joe_user_api_key_opaque_string"
	TOKEN       = "aaaaa-bbbbb-ccccc-dddd"
	EXPIRES     = "2012-04-13T13:15:00.000-05:00"
	TENANT_ID   = "12345"
	TENANT_NAME = "Opaque Name Here"

	// Example taken from http://docs.rackspace.com/auth/api/v2.0/auth-client-devguide/content/Sample_Request_Response-d1e64.html
	// and then removed extraneous records.
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
	// Same as above, but includes 'tenant' block in the Token chunk.
	SUCCESSFUL_LOGIN_RESPONSE_WITH_TENANT_IDS = `{
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
			"id": "aaaaa-bbbbb-ccccc-dddd",
			"tenant": {
				"id": "12345",
				"name": "Opaque Name Here"
			}
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
)

func TestNewIdentity(t *testing.T) {
	id := NewIdentity(USERNAME, PASSWORD, "")
	if id.Username() != USERNAME {
		t.Error("NewIdentity: expected properly set username")
		return
	}
	if id.Password() != PASSWORD {
		t.Error("NewIdentity: expected password to be set")
		return
	}
	if id.Region() != "" {
		t.Error("NewIdentity: expected no specific region to be set")
		return
	}
	if id.AuthEndpoint() != US_ENDPOINT {
		t.Error("NewIdentity: expected US endpoint used for authentication")
		return
	}
	if id.IsAuthenticated() {
		t.Error("NewIdentity: new identities are not authenticated by default")
		return
	}

	id2 := NewIdentity(USERNAME, PASSWORD, "lon")
	if id2.AuthEndpoint() != UK_ENDPOINT {
		t.Error("NewIdentity: LON region must use UK endpoint")
		return
	}
}

type testTransport struct {
	response string
	called   uint
}

func (t *testTransport) RoundTrip(req *http.Request) (rsp *http.Response, err error) {
	t.called++

	headers := make(http.Header)
	headers.Add("Content-Type", "application/xml; charset=UTF-8")

	body := ioutil.NopCloser(strings.NewReader(t.response))

	rsp = &http.Response{
		Status:           "204 OK",
		StatusCode:       200,
		Proto:            "HTTP/1.1",
		ProtoMajor:       1,
		ProtoMinor:       1,
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

func TestAuthenticationWithTenant(t *testing.T) {
	transport := &testTransport{
		response: SUCCESSFUL_LOGIN_RESPONSE_WITH_TENANT_IDS,
	}
	id := NewIdentity(USERNAME, PASSWORD, "")
	id.UseClient(&http.Client{
		Transport: transport,
	})
	err := id.Authenticate()
	if err != nil {
		t.Error("Auth:", err)
		return
	}
	if transport.called != 1 {
		t.Error("Auth: Expected HTTP request to be issued")
		return
	}
	if !id.IsAuthenticated() {
		t.Error("Auth: Expected authentication to succeed")
		return
	}
	tok, _ := id.Token()
	if tok != TOKEN {
		t.Error("Auth: Misparsed token: expected", TOKEN, "got:", tok)
		return
	}
	tenantId, _ := id.TenantId()
	if tenantId != TENANT_ID {
		t.Error("Auth: Expected tenant ID", TENANT_ID, "got:", tenantId)
		return
	}
	tenantName, _ := id.TenantName()
	if tenantName != TENANT_NAME {
		t.Error("Auth: Expected tenant name", TENANT_NAME, "got:", tenantName)
		return
	}
}

func TestAuthentication(t *testing.T) {
	transport := &testTransport{
		response: SUCCESSFUL_LOGIN_RESPONSE,
	}
	id := NewIdentity(USERNAME, PASSWORD, "")
	id.UseClient(&http.Client{
		Transport: transport,
	})
	err := id.Authenticate()
	if err != nil {
		t.Error("Auth:", err)
		return
	}
	if transport.called != 1 {
		t.Error("Auth: Expected HTTP request to be issued")
		return
	}
	if !id.IsAuthenticated() {
		t.Error("Auth: Expected authentication to succeed")
		return
	}
	tok, _ := id.Token()
	if tok != TOKEN {
		t.Error("Auth: Misparsed token: expected", TOKEN, "got:", tok)
		return
	}
	exp, _ := id.Expires()
	if exp != EXPIRES {
		t.Error("Auth: Misparsed expiration timestamp: expected", EXPIRES, "got:", exp)
		return
	}
	tenantId, _ := id.TenantId()
	if tenantId != "" {
		t.Error("Auth: unexpected tenant ID")
		return
	}
	tenantName, _ := id.TenantName()
	if tenantName != "" {
		t.Error("Auth: unexpected tenant name")
		return
	}
	svc, _ := id.ServiceCatalog()
	if len(svc) != 2 {
		t.Error("Auth: Heuristic -- service catalog count doesn't match 2; got", len(svc))
		return
	}
	roles, _ := id.Roles()
	if len(roles) != 1 {
		t.Error("Auth: Heuristic -- roles length isn't 1; got", len(roles))
		return
	}
}
