// Copyright (c) 2021 BitMaelum Authors
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

package vault

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// AccountInfo represents client account information
type AccountInfo struct {
	// Default bool             `json:"default"` // Is this the default account
	Address *address.Address `json:"address"` // The address of the account

	Name     string            `json:"name"`     // Full name of the user
	Settings map[string]string `json:"settings"` // Additional settings that can be user-defined

	// Communication and encryption information
	PrivKey   bmcrypto.PrivKey         `json:"priv_key"`        // PEM encoded private key
	PubKey    bmcrypto.PubKey          `json:"pub_key"`         // PEM encoded public key
	Pow       *proofofwork.ProofOfWork `json:"proof,omitempty"` // Proof of work
	RoutingID string                   `json:"routing_id"`      // ID of the routing used
}

// AddAccount adds a new account to the vault
func (v *Vault) AddAccount(account AccountInfo) {
	v.Store.Accounts = append(v.Store.Accounts, account)
}

// RemoveAccount removes the given account from the vault
func (v *Vault) RemoveAccount(addr address.Address) {
	k := 0
	for _, acc := range v.Store.Accounts {
		if acc.Address.String() != addr.String() {
			v.Store.Accounts[k] = acc
			k++
		}
	}
	v.Store.Accounts = v.Store.Accounts[:k]
}

var errAccountNotFound = errors.New("cannot find account")

// GetAccountInfo tries to find the given address and returns the account from the vault
func (v *Vault) GetAccountInfo(addr address.Address) (*AccountInfo, error) {
	for i := range v.Store.Accounts {
		if v.Store.Accounts[i].Address.String() == addr.String() {
			return &v.Store.Accounts[i], nil
		}
	}

	return nil, errAccountNotFound
}

// HasAccount returns true when the vault has an account for the given address
func (v *Vault) HasAccount(addr address.Address) bool {
	_, err := v.GetAccountInfo(addr)

	return err == nil
}

// // GetDefaultAccount returns the default account from the vault. This could be the one set to default, or if none found,
// // the first account in the vault. Returns nil when no accounts are present in the vault.
// func (v *Vault) GetDefaultAccount() *AccountInfo {
// 	// No accounts, return nil
// 	if len(v.Store.Accounts) == 0 {
// 		return nil
// 	}
//
// 	// Return account that is set default (the first one, if multiple)
// 	for i := range v.Store.Accounts {
// 		if v.Store.Accounts[i].Default {
// 			return &v.Store.Accounts[i]
// 		}
// 	}
//
// 	// No default found, return the first account
// 	return &v.Store.Accounts[0]
// }
