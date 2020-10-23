// Copyright (c) 2020 BitMaelum Authors
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

package internal

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// AccountInfo represents client account information
type AccountInfo struct {
	Default bool   `json:"default"` // Is this the default account
	Address string `json:"address"` // The address of the account

	Name     string            `json:"name"`     // Full name of the user
	Settings map[string]string `json:"settings"` // Additional settings that can be user-defined

	// Communication and encryption information
	PrivKey   bmcrypto.PrivKey         `json:"priv_key"`        // PEM encoded private key
	PubKey    bmcrypto.PubKey          `json:"pub_key"`         // PEM encoded public key
	Pow       *proofofwork.ProofOfWork `json:"proof,omitempty"` // Proof of work
	RoutingID string                   `json:"routing_id"`      // ID of the routing used
}

func (info *AccountInfo) AddressHash() hash.Hash {
	addr, _ := address.NewAddress(info.Address)
	return addr.Hash()
}

// OrganisationInfo represents a organisation configuration for a server
type OrganisationInfo struct {
	Addr        string                        `json:"addr"`          // org part from the bitmaelum address
	FullName    string                        `json:"name"`          // Full name of the organisation
	PrivKey     bmcrypto.PrivKey              `json:"priv_key"`      // PEM encoded private key
	PubKey      bmcrypto.PubKey               `json:"pub_key"`       // PEM encoded public key
	Pow         *proofofwork.ProofOfWork      `json:"pow,omitempty"` // Proof of work
	Validations []organisation.ValidationType `json:"validations"`   // Validations
}

// RoutingInfo represents a routing configuration for a server
type RoutingInfo struct {
	RoutingID string                   `json:"routing_id"`    // ID
	PrivKey   bmcrypto.PrivKey         `json:"priv_key"`      // PEM encoded private key
	PubKey    bmcrypto.PubKey          `json:"pub_key"`       // PEM encoded public key
	Pow       *proofofwork.ProofOfWork `json:"pow,omitempty"` // Proof of work
	Route     string                   `json:"route"`         // Route to server
}

// InfoToOrg converts organisation info to an actual organisation structure
func InfoToOrg(info OrganisationInfo) (*organisation.Organisation, error) {
	return &organisation.Organisation{
		Hash:       hash.New(info.Addr),
		FullName:   info.FullName,
		PublicKey:  info.PubKey,
		Validation: info.Validations,
	}, nil
}
