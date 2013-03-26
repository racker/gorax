package identity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	US_ENDPOINT = "https://identity.api.rackspacecloud.com/v2.0/"
	UK_ENDPOINT = "https://lon.identity.api.rackspacecloud.com/v2.0/"
)

// The Identity type encapsulates both the set of credentials used to authenticate
// against the Rackspace Identity API, as well as the relevant proof of authentication
// once acquired.
type Identity struct {
	username, apiKey, region string
	isAuthenticated          bool
	httpClient               *http.Client
	token, expires           string
	tenantId, tenantName     string
	access                   *AccessBody
}

// NewIdentity creates a new set of papers to use for authentication against the Rackspace Identity service.
// It takes a username and API key as inputs.
// Specify "" if you intend on specifying username or API key later.
// Consult with your cloud provider for your username and API key.
// The region parameter, if provided, specifies the geographical home for your account.
// Specify "" for default region (currently US).
func NewIdentity(userName, key, reg string) *Identity {
	return &Identity{
		username:   userName,
		apiKey:     key,
		region:     strings.ToUpper(reg),
		httpClient: &http.Client{},
	}
}

// SetCredentials may be used to alter the current set of credentials,
// provided the identity has not yet been authenticated.
func (id *Identity) SetCredentials(userName, key, reg string) {
	if !id.isAuthenticated {
		id.username = userName
		id.apiKey = key
		id.region = strings.ToUpper(reg)
	}
}

// Username yields the identity's user name string.
// This string is opaque to gorax.
func (id *Identity) Username() string {
	return id.username
}

// ApiKey yields the identity's API key.
// This string is opaque to gorax.
func (id *Identity) ApiKey() string {
	return id.apiKey
}

// Region yields the supplied region.
// The region returned will be in the customary all-uppercase notation.
// E.g., if you invoked NewIdentity() with a region of "lon", then this method
// will return "LON".
// If no region was set, "" is returned.
// In all other respects, this string is opaque to gorax.
func (id *Identity) Region() string {
	return id.region
}

// Token yields the authentication token.
// If not authenticated, an error is returned.
func (id *Identity) Token() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authenticated")
	}
	return id.token, nil
}

// Expires yields the token's expiration timestamp in ISO8601 format.
// If not authenticated, an error is returned.
func (id *Identity) Expires() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authentication")
	}
	return id.expires, nil
}

// TenantId yields the tenant ID.
// If not authenticated, an error is returned.
func (id *Identity) TenantId() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authenticated")
	}
	return id.tenantId, nil
}

// TenantName yields the tenant name.
// If not authenticated, an error is returned.
func (id *Identity) TenantName() (string, error) {
	if !id.IsAuthenticated() {
		return "", fmt.Errorf("Not authenticated")
	}
	return id.tenantName, nil
}

// AuthEndpoint yields which API endpoint will be used to perform the authentication.
func (id *Identity) AuthEndpoint() (ep string) {
	ep = US_ENDPOINT
	if id.region == "LON" {
		ep = UK_ENDPOINT
	}
	return
}

// IsAuthenticated reports on whether or not the credentials have been verified.
// When a new identity is created, by default it remains unauthenticated.
// Use the Authenticate() method to authenticate.
func (id *Identity) IsAuthenticated() bool {
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
	Region, TenantId, PublicURL, InternalURL string
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
func (id *Identity) ServiceCatalog() ([]CatalogEntry, error) {
	if !id.IsAuthenticated() {
		return nil, fmt.Errorf("Not authenticated")
	}
	return id.access.Access.ServiceCatalog, nil
}

// Roles yields a slice (potentially zero-length) of roles.
// An error is returned if not authenticated.
func (id *Identity) Roles() ([]Role, error) {
	if !id.IsAuthenticated() {
		return nil, fmt.Errorf("Not authenticated")
	}
	return id.access.Access.User.Roles, nil
}

// Authenticate attempts to verify this Identity object's credentials.
func (id *Identity) Authenticate() error {
	req, err := http.NewRequest("GET", id.AuthEndpoint(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := id.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Expected 200 response; got %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	jsonContainer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonContainer, &id.access)
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
func (id *Identity) UseClient(c *http.Client) {
	id.httpClient = c
}
