// vim: ts=8 sw=8 noet ai

package identity

import (
	"fmt"
	"github.com/racker/perigee"
	"net/http"
	"strings"
)

const (
	US_ENDPOINT = "https://identity.api.rackspacecloud.com/v2.0/tokens"
	UK_ENDPOINT = "https://lon.identity.api.rackspacecloud.com/v2.0/tokens"
)

// The identity (lower-case i) structure records the username, password, and
// region for the user's credentials.  In addition, it tracks whether or not
// the user is authenticated.
type identity struct {
	username, password, region string
	isAuthenticated            bool
	httpClient                 *http.Client
	token, expires             string
	tenantId, tenantName       string
	access                     *AccessBody
}

// NewIdentity creates a new set of papers to use for authentication against the Rackspace Identity service.
// It takes a username and password as inputs.
// Specify "" if you intend on specifying username or password later.
// Consult with your cloud provider for your username and password.
// The region parameter, if provided, specifies the geographical home for your account.
// Specify "" for default region (currently US).
func NewIdentity(userName, pw, reg string) *identity {
	return &identity{
		username:   userName,
		password:   pw,
		region:     strings.ToUpper(reg),
		httpClient: &http.Client{},
	}
}

// SetCredentials may be used to alter the current set of credentials,
// provided the identity has not yet been authenticated.
func (id *identity) SetCredentials(userName, pw, reg string) {
	if !id.isAuthenticated {
		id.username = userName
		id.password = pw
		id.region = strings.ToUpper(reg)
	}
}

// Username yields the identity's user name string.
// This string is opaque to gorax.
func (id *identity) Username() string {
	return id.username
}

// Password yields the identity's password.
// This string is opaque to gorax.
func (id *identity) Password() string {
	return id.password
}

// Region yields the supplied region.
// The region returned will be in the customary all-uppercase notation.
// E.g., if you invoked NewIdentity() with a region of "lon", then this method
// will return "LON".
// If no region was set, "" is returned.
// In all other respects, this string is opaque to gorax.
func (id *identity) Region() string {
	return id.region
}

// Token yields the authentication token.
// If not authenticated, an error is returned.
func (id *identity) Token() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authenticated")
	}
	return id.token, nil
}

// Expires yields the token's expiration timestamp in ISO8601 format.
// If not authenticated, an error is returned.
func (id *identity) Expires() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authentication")
	}
	return id.expires, nil
}

// TenantId yields the tenant ID.
// If not authenticated, an error is returned.
func (id *identity) TenantId() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authenticated")
	}
	return id.tenantId, nil
}

// TenantName yields the tenant name.
// If not authenticated, an error is returned.
func (id *identity) TenantName() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authenticated")
	}
	return id.tenantName, nil
}

// AuthEndpoint yields which API endpoint will be used to perform the authentication.
func (id *identity) AuthEndpoint() (ep string) {
	ep = US_ENDPOINT
	if id.region == "LON" {
		ep = UK_ENDPOINT
	}
	return
}

// IsAuthenticated reports on whether or not the credentials have been verified.
// When a new identity is created, by default it remains unauthenticated.
// Use the Authenticate() method to authenticate.
func (id *identity) IsAuthenticated() bool {
	return id.isAuthenticated
}

// AccessBody serves as a convenient wrapper for Access JSON records.
// You'll likely rarely use this type unless you intend on marshalling or unmarshalling
// Identity API JSON records yourself.
type AccessBody struct {
	Access Access
}

// Access encapsulates the API token and its relevant fields, as well as the
// services catalog that Identity API returns once authenticated.  You'll probably
// rarely use this record directly, unless you intend on marshalling or unmarshalling
// Identity API JSON records yourself.
type Access struct {
	Token          Token
	ServiceCatalog []CatalogEntry
	User           User
}

// Token encapsulates an authentication token and when it expires.  It also includes
// tenant information if available.
type Token struct {
	Id, Expires string
	Tenant      Tenant
}

// Tenant encapsulates tenant authentication information.  If, after authentication,
// no tenant information is supplied, both Id and Name will be "".
type Tenant struct {
	Id, Name string
}

// CatalogEntry encapsulates a service catalog record.
type CatalogEntry struct {
	Name, Type string
	Endpoints  []EntryEndpoint
}

// EntryEndpoint encapsulates how to get to the API of some service.
type EntryEndpoint struct {
	Region, TenantId                    string
	PublicURL, InternalURL              string
	VersionId, VersionInfo, VersionList string
}

// User encapsulates the user credentials, and provides visibility in what
// the user can do through its role assignments.
type User struct {
	Id, Name          string
	XRaxDefaultRegion string `json:"RAX-AUTH:defaultRegion"`
	Roles             []Role
}

// Role encapsulates a permission that a user can rely on.
type Role struct {
	Description, Id, Name string
}

// ServiceCatalog yields the array of services available to the user.
// An error is returned if not authenticated.
func (id *identity) ServiceCatalog() ([]CatalogEntry, error) {
	if !id.IsAuthenticated() {
		return nil, fmt.Errorf("Not authenticated")
	}
	return id.access.Access.ServiceCatalog, nil
}

// Roles yields a slice (potentially zero-length) of roles.
// An error is returned if not authenticated.
func (id *identity) Roles() ([]Role, error) {
	if !id.IsAuthenticated() {
		return nil, fmt.Errorf("Not authenticated")
	}
	return id.access.Access.User.Roles, nil
}

// These odd type definitions are required due to a known bug (more like an
// uncertainty in the semantics of the library) found in the encoding/json package.
// See http://golang.org/pkg/encoding/json/#pkg-note-BUG

type AuthContainer struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
}

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authenticate attempts to verify this Identity object's credentials.
func (id *identity) Authenticate() error {
	creds := &AuthContainer{
		Auth: Auth{
			PasswordCredentials: PasswordCredentials{
				Username: id.username,
				Password: id.password,
			},
		},
	}

	err := perigee.Post(id.AuthEndpoint(), perigee.Options{
		CustomClient: id.httpClient,
		ReqBody:      creds,
		Results:      &id.access,
	})
	if err != nil {
		return err
	}

	id.isAuthenticated = true
	id.token = id.access.Access.Token.Id
	id.expires = id.access.Access.Token.Expires
	id.tenantId = id.access.Access.Token.Tenant.Id
	id.tenantName = id.access.Access.Token.Tenant.Name
	return nil
}

// UseClient configures the identity client to use a specific net/http client.
// This allows you to configure a custom HTTP transport for specialized requirements.
// You normally wouldn't need to set this, as the net/http package makes reasonable
// choices on its own.  Customized transports are useful, however, if extra logging
// is required, or if you're using unit tests to isolate and verify correct behavior.
func (id *identity) UseClient(c *http.Client) {
	id.httpClient = c
}
