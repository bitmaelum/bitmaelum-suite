// Copyright (c) 2022 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package api

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// GetPublicKey gets public key for given address on the mail server
func (api *API) GetPublicKey(addr hash.Hash) (string, error) {
	type PubKeyOutput struct {
		PublicKey string `json:"public_key"`
	}
	output := PubKeyOutput{}

	resp, statusCode, err := api.GetJSON("/account/"+addr.String()+"/key", output)
	if err != nil {
		return "", err
	}

	if statusCode < 200 || statusCode > 299 {
		return "", GetErrorFromResponse(resp)
	}

	return output.PublicKey, nil
}

// Activate will undelete / activate an account on the resolver again
func (api *API) Activate(info vault.AccountInfo) error {
	resp, statusCode, err := api.PostJSON("/account/"+info.Address.Hash().String()+"/undelete", nil)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}

// Deactivate will delete / deactivate an account on the resolver
func (api *API) Deactivate(info vault.AccountInfo) error {
	resp, statusCode, err := api.PostJSON("/account/"+info.Address.Hash().String()+"/delete", nil)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}

// CreateAccount creates new account on server
func (api *API) CreateAccount(info vault.AccountInfo, token string) error {
	type inputCreateAccount struct {
		Addr        hash.Hash       `json:"address"`
		UserHash    hash.Hash       `json:"user_hash"`
		OrgHash     hash.Hash       `json:"org_hash"`
		Token       string          `json:"token"`
		PublicKey   bmcrypto.PubKey `json:"public_key"`
		ProofOfWork pow.ProofOfWork `json:"proof_of_work"`
	}

	input := &inputCreateAccount{
		Addr:        info.Address.Hash(),
		UserHash:    info.Address.LocalHash(),
		OrgHash:     info.Address.OrgHash(),
		Token:       token,
		PublicKey:   info.GetActiveKey().PubKey,
		ProofOfWork: *info.Pow,
	}

	resp, statusCode, err := api.PostJSON("/account", input)
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode > 299 {
		return GetErrorFromResponse(resp)
	}

	return nil
}
