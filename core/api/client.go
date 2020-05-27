package api

import (
    "crypto/tls"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/jaytaph/mailv2/core/api/types"
    "net/http"
    "strconv"
)

type Api struct {
    Host    string
    Port    int
    BaseUrl string

    client  *http.Client
}

func NewClient(host string, port int) (*Api, error) {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    api := &Api{
        Host: host,
        Port: port,
        BaseUrl : "https://" + host + ":" + strconv.Itoa(port),
        client: &http.Client{Transport: tr},
    }


    // Check info
    _, err := api.client.Get(api.BaseUrl + "/info")
    if (err != nil) {
        return nil, err
    }

    return api, nil;
}


func (api* Api) GetPublicKey(hash string) (string, error) {
    resp, err := api.client.Get(api.BaseUrl + "/account/" + hash + "/key")
    if err != nil {
        return "", err
    }

    if resp.StatusCode != http.StatusOK {
        return "", errors.New(fmt.Sprintf("Incorrect status code returned: %d", resp.StatusCode))
    }

    //b, _ := ioutil.ReadAll(resp.Body)
    //fmt.Printf("%s", string(b))
    //
    //return []byte{}, nil

    defer resp.Body.Close()

    target := types.PubKeyOutput{}
    err = json.NewDecoder(resp.Body).Decode(&target)
    if err != nil {
       return "", err
    }

    return target.PublicKey, nil
}
