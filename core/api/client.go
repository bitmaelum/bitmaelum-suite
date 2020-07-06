package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/core/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// API is a structure to connect to the server for the given account
type API struct {
	account *account.Info
	jwt     string
	client  *http.Client
}

// CreateNewClient creates a new mailserver API client
func CreateNewClient(info *account.Info) (*API, error) {
	// Create JWT token based on the private key of the user
	privKey, err := encrypt.PEMToPrivKey([]byte(info.PrivKey))
	if err != nil {
		return nil, err
	}
	hash, err := address.NewHash(info.Address)
	if err != nil {
		return nil, err
	}
	jwtToken, err := core.GenerateJWTToken(*hash, privKey)
	if err != nil {
		return nil, err
	}

	// Create API
	tr := &http.Transport{
		// Allow insecure and self-signed certificates if so configured
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Client.Server.AllowInsecure},
	}

	// If no port is present in the server, we assume port 2424
	_, _, err = net.SplitHostPort(info.Server)
	if err != nil {
		info.Server += ":2424"
	}

	if !strings.HasPrefix(info.Server, "https://") {
		info.Server = "https://" + info.Server
	}

	api := &API{
		account: info,
		jwt:     jwtToken,
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}

	return api, nil
}

// GetJSON gets JSON result from API
func (api *API) GetJSON(path string, v interface{}) error {
	body, err := api.Get(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

// Get gets raw bytes from API
func (api *API) Get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", api.account.Server+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.jwt)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("incorrect status code returned")
	}

	return ioutil.ReadAll(resp.Body)
}

// PostBytes posts to API by single bytes
func (api *API) PostBytes(path string, body []byte) error {
	return api.PostReader(path, bytes.NewBuffer(body))
}

// PostJSON posts JSON to API
func (api *API) PostJSON(path string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return api.PostReader(path, bytes.NewBuffer(b))
}

// PostReader posts to API through a reader
func (api *API) PostReader(path string, r io.Reader) error {
	req, err := http.NewRequest("POST", api.account.Server+path, r)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.jwt)

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("incorrect status code returned (%d)", resp.StatusCode)
	}

	return nil
}

// Delete from API
func (api *API) Delete(path string) error {
	req, err := http.NewRequest("DELETE", api.account.Server+path, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.jwt)

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Success codes or 404 is good
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 || resp.StatusCode == 404 {
		return nil
	}

	return fmt.Errorf("incorrect status code returned (%d)", resp.StatusCode)
}
