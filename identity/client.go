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
	"github.com/racker/gorax"
	"net/http"
)

var (
	ErrMissingCredential = &gorax.RestError{"Either a username or an apiKey must be supplied"}
)

var (
	USIdentityService = "https://identity.api.rackspacecloud.com/v2.0"
	UKIdentityService = "https://lon.identity.api.rackspacecloud.com/v2.0"
)

type authenticateWithPassword struct {
	Auth struct {
		Credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"passwordCredentials"`
	} `json:"auth"`
}

type authenticateWithAPIKey struct {
	Auth struct {
		Credentials struct {
			Username string `json:"username"`
			APIKey   string `json:"apiKey"`
		} `json:"RAX-KSKEY:apiKeyCredentials"`
	} `json:"auth"`
}

type EntryEndpoint struct {
	Region     string
	TenantId   string
	PublicURL  string
	InternaURL string
}

type CatalogEntry struct {
	Name      string
	Type      string
	Endpoints []EntryEndpoint
}

type AuthResponse struct {
	Access struct {
		Token struct {
			Id      string
			Expires string
			Tenant  struct {
				Id   string
				Name string
			}
		}
		ServiceCatalog []CatalogEntry
		User           struct {
			Id                 string
			Name               string
			ExRaxDefaultRegion string `json:"RAX-AUTH:defaultRegion"`
		}
	}
}

type KeystoneClient struct {
	username string
	password string
	apiKey   string
	client   *gorax.RestClient
}

func (k *KeystoneClient) getCredentials() (interface{}, error) {
	if len(k.apiKey) > 0 {
		data := authenticateWithAPIKey{}
		data.Auth.Credentials.Username = k.username
		data.Auth.Credentials.APIKey = k.apiKey
		return data, nil
	} else if len(k.password) > 0 {
		data := authenticateWithPassword{}
		data.Auth.Credentials.Username = k.username
		data.Auth.Credentials.Password = k.password
		return data, nil
	}

	return nil, ErrMissingCredential
}

func (k *KeystoneClient) SetDebug(debug bool) {
	k.client.SetDebug(debug)
}

func (k *KeystoneClient) Authenticate() (*AuthResponse, error) {
	creds, err := k.getCredentials()
	if err != nil {
		return nil, err
	}

	restReq := &gorax.RestRequest{
		Method: "POST",
		Path:   "/tokens",
		Body: &gorax.JSONRequestBody{
			Object: creds,
		},
		ExpectedStatusCodes: []int{http.StatusOK},
	}

	resp, err := k.client.PerformRequest(restReq)

	if err != nil {
		return nil, err
	}

	authResponse := &AuthResponse{}
	err = resp.DeserializeBody(authResponse)

	return authResponse, err
}

func MakePasswordKeystoneClient(url string, username string, password string) *KeystoneClient {
	return &KeystoneClient{
		client:   gorax.MakeRestClient(url),
		username: username,
		password: password,
	}
}

func MakeAPIKeyKeystoneClient(url string, username string, apiKey string) *KeystoneClient {
	return &KeystoneClient{
		client:   gorax.MakeRestClient(url),
		username: username,
		apiKey:   apiKey,
	}
}
