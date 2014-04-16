/*
Copyright 2013 Rackspace

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS-IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package identity

import (
	"fmt"
	"sync"
	"time"

	"github.com/racker/gorax"
)

var (
	USExpiresFormat = "2006-01-02T15:04:05.000-07:00"
	UKExpiresFormat = "2006-01-02T15:04:05.000Z"
	ExpireDelta     = time.Duration(5) * time.Minute
)

// The KeystoneAuthMiddleware object is a filter through which all requests flow,
// ensuring that the principal making requests against a Cloud Monitoring account actually
// has the required privileges to do so.
type KeystoneAuthMiddleware struct {
	tenantId       string
	token          string
	expires        time.Time
	keystoneClient *KeystoneClient
	refreshLock    sync.Mutex
}

// MakeKeystonePasswordMiddleware creates a middleware request object to the API to use the Keystone authentication interface.
// A side-effect of this function is the creation of a Keystone client interface object in debug-mode.
// This procedure assumes username/password authentication.
func MakeKeystonePasswordMiddleware(region string, username string, password string) *KeystoneAuthMiddleware {
	m := &KeystoneAuthMiddleware{
		keystoneClient: MakePasswordKeystoneClient(region, username, password),
		expires:        time.Time{},
		refreshLock:    sync.Mutex{},
	}
	m.keystoneClient.SetDebug(false)
	return m
}

// MakeKeystoneAPIKeyMiddleware creates a middleware request object to the API to use the Keystone authentication interface.
// A side-effect of this function is the creation of a Keystone client interface object in debug-mode.
// This procedure assumes you already have a valid API key for the principal making the requests.
func MakeKeystoneAPIKeyMiddleware(region string, username string, apiKey string) *KeystoneAuthMiddleware {
	m := &KeystoneAuthMiddleware{
		keystoneClient: MakeAPIKeyKeystoneClient(region, username, apiKey),
		expires:        time.Time{},
		refreshLock:    sync.Mutex{},
	}
	m.keystoneClient.SetDebug(false)
	return m
}

// This HandleRequest method performs user authentication against a Keystone REST API.
//
// If the request has timed out (e.g., as by exceeding its expiry timeout), it returns an error out of hand.  No attempt to use REST resources occurs.
// Otherwise, if the username and password work for the principal, and the request hasn't timed out, a new expiry timestamp is calculated, thus extending
// the life of the request.  Additionally, the request receives an X-Auth-Token MIME header appropriate for the principal making the request, and the
// resource path is "redirected" to that appropriate for the principal's unique set of resources.
func (m *KeystoneAuthMiddleware) HandleRequest(req *gorax.RestRequest) (*gorax.RestRequest, error) {
	m.refreshLock.Lock()
	defer m.refreshLock.Unlock()

	if time.Now().Add(ExpireDelta).After(m.expires) {
		result, err := m.keystoneClient.Authenticate()
		if err != nil {
			return nil, err
		}

		var expires time.Time

		expires, err = time.Parse(USExpiresFormat, result.Access.Token.Expires)

		if err != nil {
			expires, err = time.Parse(UKExpiresFormat, result.Access.Token.Expires)
			if err != nil {
				return nil, fmt.Errorf("unable to parse token expiration time: %s", result.Access.Token.Expires)
			}
		}

		m.tenantId = result.Access.Token.Tenant.Id
		m.token = result.Access.Token.Id
		m.expires = expires
	}

	req.Header.Set("X-Auth-Token", m.token)
	req.Path = "/" + m.tenantId + req.Path

	return req, nil
}
