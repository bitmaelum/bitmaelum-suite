package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/ticket"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/ernesto-jimenez/httplogger"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
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
		transport = httplogger.NewLoggedTransport(transport, newLogger())
	}

	api := &API{
		host: canonicalHost(opts.Host),
		client: &http.Client{
			Transport: transport,
			Timeout: 15 * time.Second,
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
		transport = httplogger.NewLoggedTransport(transport, newLogger())
	}

	if info != nil {
		// Create JWT token based on the private key of the user
		hash, err := address.NewHash(info.Address)
		if err != nil {
			return nil, err
		}
		jwtToken, err = internal.GenerateJWTToken(*hash, info.PrivKey)
		if err != nil {
			return nil, err
		}
	}

	api := &API{
		host: canonicalHost(opts.Host),
		client: &http.Client{
			Transport: transport,
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

	err = json.Unmarshal(body, v)
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

type httpLogger struct {
}

func newLogger() *httpLogger {
	logrus.SetLevel(logrus.TraceLevel)
	return &httpLogger{}
}

func (l *httpLogger) LogRequest(req *http.Request) {
	var err error
	save := req.Body
	if req.Body != nil {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return
		}
	}

	logrus.Tracef(
		"Request %s %s",
		req.Method,
		req.URL.String(),
	)

	for k, v := range req.Header {
		logrus.Tracef("HEADER: %s : %s\n", k, v)
	}

	var b bytes.Buffer
	if req.Body != nil {
		var dest io.Writer = &b
		_, err = io.Copy(dest, req.Body)
	}

	req.Body = save

	logrus.Tracef("log: %s", b.String())
}

func (l *httpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		logrus.Trace("log: ", err)
	} else {
		logrus.Tracef(
			"Response method=%s status=%d durationMs=%d %s",
			req.Method,
			res.StatusCode,
			duration,
			req.URL.String(),
		)
		for k, v := range res.Header {
			logrus.Tracef("HEADER: %s : %s\n", k, v)
		}

	}
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
