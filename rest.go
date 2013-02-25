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

package gorax

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/http/httputil"
)

type RestError struct {
	ErrorString string
}

func (e *RestError) Error() string {
	return e.ErrorString
}

type RestRequest struct {
	Method              string
	Path                string
	Header              http.Header
	Body                RequestBody
	ExpectedStatusCodes []int
}

type RequestBody interface {
	ContentType() (string, error)
	ContentLength() (int64, error)
	Body() (io.Reader, error)
}

type JSONRequestBody struct {
	Object interface{}
	data   []byte
}

func (b *JSONRequestBody) marshal() error {
	if b.data != nil {
		return nil
	}

	data, err := json.Marshal(b.Object)
	if err == nil {
		b.data = data
	}

	return err
}

func (b *JSONRequestBody) ContentType() (string, error) {
	return "application/json", nil
}

func (b *JSONRequestBody) ContentLength() (int64, error) {
	err := b.marshal()
	if err != nil {
		return 0, err
	}

	return int64(len(b.data)), nil
}

func (b *JSONRequestBody) Body() (io.Reader, error) {
	err := b.marshal()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b.data), nil
}

type RestResponse struct {
	*http.Response
}

func (r *RestResponse) DeserializeBody(target interface{}) error {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	}

	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))

	if err != nil {
		return err
	}

	switch contentType {
	case "application/json":
		return json.Unmarshal(data, target)
	}

	return fmt.Errorf("unsupported Content-Type: %s", r.Header.Get("Content-Type"))
}

type RequestMiddleware interface {
	HandleRequest(*RestRequest) (*RestRequest, error)
}

type RestClient struct {
	BaseUrl            string
	RequestMiddlewares []RequestMiddleware
	Debug              bool
	client             *http.Client
}

func MakeRestClient(url string) *RestClient {
	return &RestClient{
		BaseUrl:            url,
		RequestMiddlewares: []RequestMiddleware{},
		Debug:              false,
		client:             &http.Client{},
	}
}

func (c *RestClient) PerformRequest(restReq *RestRequest) (*RestResponse, error) {
	var err error
	var body io.Reader

	// Request Middlewares shouldn't have to worry about a nil Header
	if restReq.Header == nil {
		restReq.Header = http.Header{}
	}

	for _, middleware := range c.RequestMiddlewares {
		restReq, err = middleware.HandleRequest(restReq)

		if err != nil {
			return nil, err
		}
	}

	if restReq.Body != nil {
		body, err = restReq.Body.Body()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(restReq.Method, c.BaseUrl+restReq.Path, body)
	req.Header = restReq.Header

	if len(req.Header.Get("Accept")) == 0 {
		req.Header.Set("Accept", "application/json")
	}

	if restReq.Body != nil {
		if err != nil {
			return nil, err
		}

		contentLength, err := restReq.Body.ContentLength()
		if err != nil {
			return nil, err
		}

		req.ContentLength = contentLength

		contentType, err := restReq.Body.ContentType()
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", contentType)
	}

	if c.Debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("")
		fmt.Println(string(dump))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("")
		fmt.Println(string(dump))
	}

	if restReq.ExpectedStatusCodes != nil {
		for _, value := range restReq.ExpectedStatusCodes {
			if value == resp.StatusCode {
				return &RestResponse{resp}, nil
			}
		}

		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}

	return &RestResponse{resp}, nil
}

func (c *RestClient) SetDebug(debug bool) {
	c.Debug = debug
}
