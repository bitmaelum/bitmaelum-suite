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
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/ernesto-jimenez/httplogger"
)

type jsonOut map[string]interface{}

var errNoSuccess = errors.New("operation was not successful")

// API is a structure to connect to the server for the given account
type API struct {
	client *http.Client   // internal HTTP client
	host   string         // host to the server
	jwt    string         // optional JWT token (for client's authorized communication)
	ticket *ticket.Ticket // optional ticket for communication
}

// ClientOpts allows you to configure the API client
type ClientOpts struct {
	Host          string
	AllowInsecure bool
	Debug         bool
}

// NewAnonymous creates a new client that connects anonymously to a BitMaelum server (normally
// server-to-server communication)
func NewAnonymous(opts ClientOpts) (*API, error) {
	var transport http.RoundTripper
	transport = &http.Transport{
		// Allow insecure and self-signed certificates if so configured
		TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.AllowInsecure},
	}

	if opts.Debug {
		// Wrap transport in debug logging
		transport = httplogger.NewLoggedTransport(transport, internal.NewHTTPLogger())
	}

	api := &API{
		host: CanonicalHost(opts.Host),
		client: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
		},
	}

	return api, nil
}

// NewAuthenticated creates a new client that connects with the specified account to a BitMaelum server (normally
// client-to-server communication)
func NewAuthenticated(info *internal.AccountInfo, opts ClientOpts) (*API, error) {
	var jwtToken string
	var transport http.RoundTripper

	transport = &http.Transport{
		// Allow insecure and self-signed certificates if so configured
		TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.AllowInsecure},
	}

	if opts.Debug {
		// Wrap transport in debug logging
		transport = httplogger.NewLoggedTransport(transport, internal.NewHTTPLogger())
	}

	if info != nil {
		// Create JWT token based on the private key of the user
		addr, _ := address.NewAddress(info.Address)
		addrHash := addr.Hash()
		var err error
		jwtToken, err = internal.GenerateJWTToken(addrHash, info.PrivKey)
		if err != nil {
			return nil, err
		}
	}

	api := &API{
		host: CanonicalHost(opts.Host),
		client: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		jwt: jwtToken,
	}

	return api, nil
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
func (api *API) GetReader(path string) (r io.Reader, statusCode int, err error) {
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
	if api.jwt != "" {
		req.Header.Set("Authorization", "Bearer "+api.jwt)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, 0, err
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

func getErrorFromResponse(body []byte) error {
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
