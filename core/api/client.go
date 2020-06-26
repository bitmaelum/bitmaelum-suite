package api

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/bitmaelum/bitmaelum-server/core/encrypt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
)

type Api struct {
    account     *core.AccountInfo
    jwt         string
    client      *http.Client
}

// Create a new mailserver API client
func CreateNewClient(ai *core.AccountInfo) (*Api, error) {
    // Create JWT token based on the private key of the user
    privKey, err := encrypt.PEMToPrivKey([]byte(ai.PrivKey))
    if err != nil {
        return nil, err
    }
    jwtToken, err := core.GenerateJWTToken(core.StringToHash(ai.Address), privKey)
    if err != nil {
        return nil, err
    }

    // Create API
    tr := &http.Transport{
        // Allow insecure and self-signed certificates if so configured
        TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Client.Server.AllowInsecure},
    }

    // If no port is present in the server, we assume port 2424
    if ! strings.Contains(ai.Server, ":") {
        ai.Server = ai.Server + ":2424"
    }
    if ! strings.HasPrefix(ai.Server, "https://") {
        ai.Server = "https://" + ai.Server
    }

    api := &Api{
        account: ai,
        jwt: jwtToken,
        client: &http.Client{
            Transport: tr,
            Timeout: 30 * time.Second,
        },
    }

	return api, nil
}

// Get JSON result from API
func (api *Api) GetJSON(path string, v interface{}) error {
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

// Get raw bytes from API
func (api *Api) Get(path string) ([]byte, error) {
    req, err := http.NewRequest("GET", api.account.Server + path, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + api.jwt)

    resp, err := api.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
        return nil, errors.New("incorrect status code returned")
    }

    return ioutil.ReadAll(resp.Body)
}

// Post to API
func (api *Api) Post(path string, body interface{}) error {
    bodyBytes, err := json.Marshal(body)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", api.account.Server + path, bytes.NewBuffer(bodyBytes))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + api.jwt)

    resp, err := api.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
        return errors.New(fmt.Sprintf("incorrect status code returned (%d)", resp.StatusCode))
    }

    return nil
}
