package resolve

import (
    "bytes"
    "encoding/json"
    "errors"
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

func NewRemoteRepository(baseUrl string) Repository {
    return &remoteRepo{
        BaseUrl: baseUrl,
        client: &http.Client{},
    }
}

func (r *remoteRepo) Retrieve(hash string) (*ResolveInfo, error) {
    response, err := r.client.Get(r.BaseUrl + "/" + hash)
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

func (r *remoteRepo) Upload(hash, pubKey, address, signature string) error {
    data := &KeyUpload{
        PublicKey: pubKey,
        Address:   address,
        Signature: signature,
    }

    byteBuf, err := json.Marshal(&data)
    if err != nil {
        return err
    }

    response, err := r.client.Post(r.BaseUrl + "/" + hash, "application/json", bytes.NewBuffer(byteBuf))
    if err != nil {
        return err
    }

    if (response.StatusCode >= 200 && response.StatusCode <= 299) {
        return nil
    }

    return errors.New("error status returned from account server")
}
