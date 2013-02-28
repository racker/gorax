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

// The RestRequest object encapsulates a single RESTful request.
//
// The Method field indicates what operation to perform on the web server.
// Unless you know what you're doing, you should restrain your methods to the basic set of GET, PUT, POST, and DELETE.
//
// Path indicates the specific resource the request is being applied against.
// When making a RestClient object, you'll be able to specify a prefix which scopes this field.
// Thus, the Path field is always relative to the RestClient's base URL.
// See also the RestClient type.
//
// Header provides a mapping of MIME header keys to one or more values for each key.
// Refer to http://godoc.org/net/http#Header for more details.
// This field may be set to nil if the client needs no headers.
// However, upon making the first REST request, a nil Header field will be initialized to an empty Header map.
//
// Body provides an optional body for the request.
//
// ExpectedStatusCodes provides a set of response codes considered to be valid for the request.
// This field may also be nil if you just don't care about response checking.
type RestRequest struct {
	Method              string
	Path                string
	Header              http.Header
	Body                RequestBody
	ExpectedStatusCodes []int
}

// The RequestBody interface represents an abstract concept of a hunk of data passed along with an HTTP request.
// Bodies may be part of a request, a response, or both, depending upon the resource accessed and the method used.
//
// The ContentType() (string, error) method yields the official, even if experimental, MIME type string for the content in the body.
// The ContentLength() (int64, error) method yields its length.  Finally, the Body() (io.Reader, error) method yields a reader interface
// that allows the client software access to the contents of the data.
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

// DeserializeBody() recreates an object graph as denoted in the REST response's body.
// 
// This method will return an error if the response's content-type cannot be recognized.
// Presently, the RestResponse type only supports application/json; future versions of this
// type may support additional types.  Refer to http://godoc.org/encoding/json for more information.
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

// The RestClient object encapsulates a connection to a RESTful service.  Several assumptions are
// made about this service:
//
// 1.  All resources exposed by the service sits under a single root path, referenced by the field named BaseUrl.
// 2.  All requests made to the service pass through a set of "middleware" filters (more generally, "middlewares").
// This set may be empty, and may be inspected via the RequestMiddlewares field.
// 3.  If any one of the filters produces an error, the entire request fails.
// 4.  Filters may alter the original request.  For example, an authentication filter may inject an authentication token,
// while a tracing filter may inject a unique tracing token for logging purposes.
//
// Additionally, the Debug field indicates whether or not request/response logging (usually to stdout) occurs.
type RestClient struct {
	BaseUrl            string
	RequestMiddlewares []RequestMiddleware
	Debug              bool
	client             *http.Client
}

// MakeRestClient() creates a new RestClient reference to a RESTful service.  The provided URL sets the BaseUrl of
// the client, which scopes the resource paths of all RestRequests used to invoke services.  This function can never
// fail except in out-of-memory situations.
func MakeRestClient(url string) *RestClient {
	return &RestClient{
		BaseUrl:            url,
		RequestMiddlewares: []RequestMiddleware{},
		Debug:              false,
		client:             &http.Client{},
	}
}

// PerformRequest() returns the RESTful response to a RESTful request against a web server.
// The request must be encapsulated in a RestRequest object.
//
// This function will return with an error if the REST server provides a response which is not anticipated by the client software.
// To configure the list of anticipated response codes, the RestRequest must have a non-nil ExpectedStatusCodes field value.
// See the RestRequest type for more details.
// 
// If the request is in debug mode, diagnostic dumps of the HTTP traffic will appear on stdout.
//
// PerformRequest() computes the content length and type from the request's configured body, if any exists.
//
// If the accepted response type isn't explicitly set prior to calling PerformRequest(), it will be set to application/json by default.
//
// PerformRequest() takes care of applying any transformations to your request via the client's configured set of middleware filters.
// For example, if a client has a password authenticator middleware filter configured for it, then all requests made through that client
// will have its provided username and password fields authenticated in prior to the actual and intended request handler for the specified
// resource.
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

// The SetDebug() function either enables (true) or disables (false) diagnostic output to stdout
// of all traffic (request and response alike) made through the client.
func (c *RestClient) SetDebug(debug bool) {
	c.Debug = debug
}
