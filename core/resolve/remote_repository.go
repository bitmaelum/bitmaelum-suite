package resolve

import (
    "bytes"
    "encoding/json"
    "errors"
    "github.com/bitmaelum/bitmaelum-server/core"
    "io/ioutil"
    "net/http"
)

type remoteRepo struct {
    BaseUrl     string
    client      *http.Client
}

type KeyUpload struct {
    PublicKey   string  `json:"public_key"`
    Address     string  `json:"address"`
    Signature   string  `json:"signature"`
}

type KeyDownload struct {
    Hash        string  `json:"hash"`
    PublicKey   string  `json:"public_key"`
    Address     string  `json:"address"`
}

// Create new remote resolve repository
func NewRemoteRepository(baseUrl string) Repository {
    return &remoteRepo{
        BaseUrl: baseUrl,
        client: &http.Client{},
    }
}

// Resolve
func (r *remoteRepo) Resolve(addr core.HashAddress) (*ResolveInfo, error) {
    response, err := r.client.Get(r.BaseUrl + "/" + addr.String())
    if err != nil {
        return nil, errors.New("Error while retrieving key")
    }

    if response.StatusCode == 404 {
        return nil, errors.New("Key not found")
    }

    if response.StatusCode == 200 {
        res, err := ioutil.ReadAll(response.Body)
        if err != nil {
            return nil, errors.New("Error while retrieving key")
        }

        ri := &ResolveInfo{}
        err = json.Unmarshal(res, &ri)
        if err != nil {
            return nil, errors.New("Error while retrieving key")
        }

        return ri, nil
    }

    return nil, errors.New("Error while retrieving key")
}

func (r *remoteRepo) Upload(addr core.HashAddress, pubKey, address, signature string) error {
    data := &KeyUpload{
        PublicKey: pubKey,
        Address:   address,
        Signature: signature,
    }

    byteBuf, err := json.Marshal(&data)
    if err != nil {
        return err
    }

    response, err := r.client.Post(r.BaseUrl + "/" + addr.String(), "application/json", bytes.NewBuffer(byteBuf))
    if err != nil {
        return err
    }

    body, err := ioutil.ReadAll(response.Body)
    if (response.StatusCode >= 200 && response.StatusCode <= 299) {
        return nil
    }

    return errors.New(string(body))
}
