package api

import (
	"github.com/bitmaelum/bitmaelum-server/core"
)

// GetPublicKey gets public key for given address on the mail server
func (api *API) GetPublicKey(addr core.HashAddress) (string, error) {
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
func (api *API) CreateAccount(ai core.AccountInfo, token string) error {
	type InputCreateAccount struct {
		Addr        core.HashAddress `json:"address"`
		Token       string           `json:"token"`
		PublicKey   string           `json:"public_key"`
		ProofOfWork struct {
			Bits  int    `json:"bits"`
			Proof uint64 `json:"proof"`
		} `json:"proof_of_work"`
	}

	addr, _ := core.NewAddressFromString(ai.Address)

	input := &InputCreateAccount{
		Addr:      addr.Hash(),
		Token:     token,
		PublicKey: ai.PubKey,
		ProofOfWork: struct {
			Bits  int    `json:"bits"`
			Proof uint64 `json:"proof"`
		}{
			Bits:  ai.Pow.Bits,
			Proof: ai.Pow.Proof,
		},
	}
	return api.PostJSON("/account", input)
}
