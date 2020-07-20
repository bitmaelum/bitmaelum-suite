package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type jsonOut map[string]interface{}

var errNoSuccess = errors.New("operating was not successful")

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
}

// NewAnonymous creates a new client that connects anonymously to a BitMaelum server (normally
// server-to-server communication)
func NewAnonymous(opts ClientOpts) (*API, error) {
	api := &API{
		host: canonicalHost(opts.Host),
		client: &http.Client{
			Transport: &http.Transport{
				// Allow insecure and self-signed certificates if so configured
				TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.AllowInsecure},
			},
			Timeout: 30 * time.Second,
		},
	}

	return api, nil
}

// NewAuthenticated creates a new client that connects with the specified account to a BitMaelum server (normally
// client-to-server communication)
func NewAuthenticated(info *account.Info, opts ClientOpts) (*API, error) {
	var jwtToken string

	if info != nil {
		// Create JWT token based on the private key of the user
		privKey, err := encrypt.PEMToPrivKey([]byte(info.PrivKey))
		if err != nil {
			return nil, err
		}
		hash, err := address.NewHash(info.Address)
		if err != nil {
			return nil, err
		}
		jwtToken, err = internal.GenerateJWTToken(*hash, privKey)
		if err != nil {
			return nil, err
		}
	}

	api := &API{
		host: canonicalHost(opts.Host),
		client: &http.Client{
			Transport: &http.Transport{
				// Allow insecure and self-signed certificates if so configured
				TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.AllowInsecure},
			},
			Timeout: 30 * time.Second,
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

	return api.do(req)
}

// GetJSON gets JSON result from API
func (api *API) GetJSON(path string, v interface{}) (status int, err error) {
	body, status, err := api.Get(path)
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return status, err
	}

	return status, nil
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

	return api.do(req)
}

// Delete from API
func (api *API) Delete(path string) (body []byte, statusCode int, err error) {
	req, err := http.NewRequest("DELETE", api.host+path, nil)
	if err != nil {
		return nil, 0, err
	}

	return api.do(req)
}

// setTicketHeader sets a ticket so api requests add them to the header when doing requests
func (api *API) setTicketHeader(t ticket.Ticket) {
	api.ticket = &t
}

// do actually does the given request
func (api *API) do(req *http.Request) (body []byte, statusCode int, err error) {
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
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err = ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

// canonicalHost returns a given host in the form of http(s)://<host>:<port>
func canonicalHost(host string) string {
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
