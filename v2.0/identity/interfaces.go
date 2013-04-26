// vim: ts=8 sw=8 noet ai

package identity

// The Identity interface encapsulates both the set of credentials used to
// authenticate against the Rackspace Identity API, as well as the relevant
// proof of authentication once acquired.
type Identity interface {
	SetCredentials(userName, password, reg string)
	Username() string
	Password() string
	Region() string
	Token() (string, error)
	Expires() (string, error)
	TenantId() (string, error)
	TenantName() (string, error)
	AuthEndpoint() (ep string)
	IsAuthenticated() bool
	ServiceCatalog() ([]CatalogEntry, error)
	Roles() ([]Role, error)
	Authenticate() error
}
