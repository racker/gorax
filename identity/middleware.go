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
	"github.com/racker/gorax"
	"sync"
	"time"
)

var (
	USExpiresFormat = "2006-01-02T15:04:05.000-07:00"
	UKExpiresFormat = "2006-01-02T15:04:05.000Z"
	ExpireDelta     = time.Duration(5) * time.Minute
)

type KeystoneAuthMiddleware struct {
	tenantId       string
	token          string
	expires        time.Time
	keystoneClient *KeystoneClient
	refreshLock    sync.Mutex
}

func MakeKeystonePasswordMiddleware(region string, username string, password string) *KeystoneAuthMiddleware {
	m := &KeystoneAuthMiddleware{
		keystoneClient: MakePasswordKeystoneClient(region, username, password),
		expires:        time.Time{},
		refreshLock:    sync.Mutex{},
	}
	m.keystoneClient.SetDebug(true)
	return m
}

func MakeKeystoneAPIKeyMiddleware(region string, username string, apiKey string) *KeystoneAuthMiddleware {
	m := &KeystoneAuthMiddleware{
		keystoneClient: MakeAPIKeyKeystoneClient(region, username, apiKey),
		expires:        time.Time{},
		refreshLock:    sync.Mutex{},
	}
	m.keystoneClient.SetDebug(true)
	return m
}

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
