// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/ernesto-jimenez/httplogger"
)

type jsonOut map[string]interface{}

var errNoSuccess = errors.New("operation was not successful")

type serverErrorFunc func(*http.Request, *http.Response)

// API is a structure to connect to the server for the given account
type API struct {
	client    *http.Client      // internal HTTP client
	host      string            // host to the server
	address   address.Address   // optional address of the user to generate the jwt token
	key       *bmcrypto.PrivKey // optional private key of the user to generate the jwt token
	ticket    *ticket.Ticket    // optional ticket for communication
	errorFunc serverErrorFunc   // optional function to call when a server call fails
}

// ClientOpts allows you to configure the API client
type ClientOpts struct {
	Host          string
	AllowInsecure bool
	Address       address.Address
	Key           *bmcrypto.PrivKey
	Debug         bool
	ErrorFunc     serverErrorFunc
}

// NewAnonymous creates a new client that connects anonymously to a BitMaelum server (normally server-to-server communications)
func NewAnonymous(host string, f serverErrorFunc) (*API, error) {
	return NewClient(ClientOpts{
		Host:          host,
		AllowInsecure: config.Client.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
		ErrorFunc:     f,
	})
}

// NewAuthenticated creates a new client that connects to a BitMaelum Server with specific credentials (normally client-to-server communications)
func NewAuthenticated(addr address.Address, key *bmcrypto.PrivKey, host string, f serverErrorFunc) (*API, error) {
	return NewClient(ClientOpts{
		Host:          host,
		Address:       addr,
		Key:           key,
		AllowInsecure: config.Client.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
		ErrorFunc:     f,
	})
}

// NewClient creates a new client based on the given options
func NewClient(opts ClientOpts) (*API, error) {
	var transport http.RoundTripper

	transport = &http.Transport{
		// Allow insecure and self-signed certificates if so configured
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.AllowInsecure,
		},
	}

	if opts.Debug {
		// Wrap transport in debug logging
		transport = httplogger.NewLoggedTransport(transport, internal.NewHTTPLogger())
	}

	return &API{
		host: CanonicalHost(opts.Host),
		client: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
		},
		address:   opts.Address,
		key:       opts.Key,
		errorFunc: opts.ErrorFunc,
	}, nil
}

// Get gets raw bytes from API
func (api *API) Get(path string) (body []byte, statusCode int, err error) {
	req, err := http.NewRequest("GET", api.host+path, nil)
	if err != nil {
		return nil, 0, err
	}

	r, statusCode, err := api.do(req)
	if err != nil {
		return nil, statusCode, err
	}
	body, err = ioutil.ReadAll(r)
	_ = r.Close()
	return body, statusCode, err
}

// GetReader returns a io.Reader from the API. It is the same as api.Get(), except this one already has read the whole
// body.
func (api *API) GetReader(path string) (r io.ReadCloser, statusCode int, err error) {
	req, err := http.NewRequest("GET", api.host+path, nil)
	if err != nil {
		return nil, 0, err
	}

	return api.do(req)
}

// GetJSON gets JSON result from API
func (api *API) GetJSON(path string, v interface{}) (body []byte, status int, err error) {
	body, status, err = api.Get(path)
	if err != nil {
		return body, status, err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return body, status, err
	}

	return body, status, nil
}

// Post posts to API by single bytes
func (api *API) Post(path string, data []byte) (body []byte, statusCode int, err error) {
	return api.PostReader(path, bytes.NewBuffer(data))
}

// PostJSON posts JSON to API
func (api *API) PostJSON(path string, data interface{}) (body []byte, statusCode int, err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	return api.PostReader(path, bytes.NewBuffer(b))
}

// PostReader posts to API through a reader
func (api *API) PostReader(path string, r io.Reader) (body []byte, statusCode int, err error) {
	req, err := http.NewRequest("POST", api.host+path, r)
	if err != nil {
		return nil, 0, err
	}

	r, statusCode, err = api.do(req)
	if err != nil {
		return nil, statusCode, err
	}
	body, err = ioutil.ReadAll(r)
	return body, statusCode, err
}

// Delete from API
func (api *API) Delete(path string) (body []byte, statusCode int, err error) {
	req, err := http.NewRequest("DELETE", api.host+path, nil)
	if err != nil {
		return nil, 0, err
	}

	r, statusCode, err := api.do(req)
	if err != nil {
		return nil, statusCode, err
	}
	body, err = ioutil.ReadAll(r)
	return body, statusCode, err
}

// setTicketHeader sets a ticket so api requests add them to the header when doing requests
func (api *API) setTicketHeader(t ticket.Ticket) {
	api.ticket = &t
}

// do does the actual request and returns body reader, status code and error
func (api *API) do(req *http.Request) (body io.ReadCloser, statusCode int, err error) {
	req.Header.Set("Content-Type", "application/json")
	if api.ticket != nil {
		req.Header.Set(ticket.TicketHeader, api.ticket.ID)
	}
	if api.key != nil {
		jwtToken, err := GenerateJWTToken(api.address.Hash(), *api.key)
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	// Call the error function if there was an error
	if resp.StatusCode >= 400 && api.errorFunc != nil {
		api.errorFunc(req, resp)
	}

	return resp.Body, resp.StatusCode, nil
}

// CanonicalHost returns a given host in the form of http(s)://<host>:<port>
func CanonicalHost(host string) string {
	// If no port is present in the server, we assume port 2424
	_, _, err := net.SplitHostPort(host)
	if err != nil {
		host += ":2424"
	}

	// if no protocol is given, assume https://
	if !strings.Contains(host, "://") {
		host = "https://" + host
	}

	return host
}

// GetErrorFromResponse will return an error generated from the body
func GetErrorFromResponse(body []byte) error {
	type errorStatus struct {
		Error  bool   `json:"error"`
		Status string `json:"status"`
	}

	s := &errorStatus{}
	err := json.Unmarshal(body, &s)
	if err != nil {
		return errNoSuccess
	}

	return errors.New(s.Status)
}

func isErrorResponse(body []byte) bool {
	type errorStatus struct {
		Error  bool   `json:"error"`
		Status string `json:"status"`
	}

	s := &errorStatus{}
	err := json.Unmarshal(body, &s)
	if err != nil {
		return false
	}

	return s.Error
}
