package api

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// GetPublicKey gets public key for given address on the mail server
func (api *API) GetPublicKey(addr address.HashAddress) (string, error) {
	type PubKeyOutput struct {
		PublicKey string `json:"public_key"`
	}
	output := PubKeyOutput{}

	err := api.GetJSON("/account/"+addr.String()+"/key", output)
	if err != nil {
		return "", err
	}

	return output.PublicKey, nil
}

// CreateAccount creates new account on server
func (api *API) CreateAccount(info account.Info, token string) error {
	type InputCreateAccount struct {
		Addr        address.HashAddress `json:"address"`
		Token       string              `json:"token"`
		PublicKey   string              `json:"public_key"`
		ProofOfWork pow.ProofOfWork     `json:"proof_of_work"`
	}

	addr, _ := address.New(info.Address)

	input := &InputCreateAccount{
		Addr:        addr.Hash(),
		Token:       token,
		PublicKey:   info.PubKey,
		ProofOfWork: info.Pow,
	}
	return api.PostJSON("/account", input)
}
