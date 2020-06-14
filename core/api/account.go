package api

import (
    "github.com/bitmaelum/bitmaelum-server/core"
)


// Get public key for given address on the mail server
func (api* Api) GetPublicKey(addr core.HashAddress) (string, error) {
    type PubKeyOutput struct {
        PublicKey   string  `json:"public_key"`
    }
    output := PubKeyOutput{}

    err := api.GetJSON("/account/" + addr.String() + "/key", output)
    if err != nil {
        return "", err
    }

    return output.PublicKey, nil
}

// Create new account on server
func (api *Api) CreateAccount(addr core.HashAddress, token string) error {
    err := api.Post("/account", nil)

    return err
}