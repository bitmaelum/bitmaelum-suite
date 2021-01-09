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

package config

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/spf13/afero"
)

// RoutingConfig holds routing configuration for the mail server
type RoutingConfig struct {
	Version   int               `json:"version"`
	RoutingID string            `json:"routing_id"`
	KeyPair   *bmcrypto.KeyPair `json:"keypair,omitempty"`
}

type routingConfigV1 struct {
	RoutingID  string            `json:"routing_id"`
	PrivateKey *bmcrypto.PrivKey `json:"private_key,omitempty"`
	PublicKey  *bmcrypto.PubKey  `json:"public_key,omitempty"`
}

// Routing keeps the routing ID and keys
var Routing RoutingConfig

// ReadRouting will read the routing file and merge it into the server configuration
func ReadRouting(p string) error {
	f, err := fs.Open(p)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = readConfigV0(p, data)
	if err != nil {
		err = readConfigV1(p, data)
	}

	return err
}

func readConfigV0(p string, data []byte) error {
	tmp := routingConfigV1{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	// Migrate to new structure
	Routing.KeyPair = &bmcrypto.KeyPair{
		Generator:   "",
		FingerPrint: tmp.PublicKey.Fingerprint(),
		PrivKey:     *tmp.PrivateKey,
		PubKey:      *tmp.PublicKey,
	}

	return SaveRouting(p, &Routing)
}

func readConfigV1(p string, data []byte) error {
	tmp := routingConfigV1{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	Routing = RoutingConfig{}
	return json.Unmarshal(data, &Routing)
}

// SaveRouting will save the routing into a file. It will overwrite if exists
func SaveRouting(p string, routing *RoutingConfig) error {
	data, err := json.MarshalIndent(routing, "", "  ")
	if err != nil {
		return err
	}

	err = fs.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, p, data, 0600)
}

// GenerateRoutingFromMnemonic generates a new routing file from the given seed
func GenerateRoutingFromMnemonic(mnemonic string) (*RoutingConfig, error) {
	kp, err := internal.GenerateKeypairFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	id := hex.EncodeToString(kp.PubKey.K.(ed25519.PublicKey))
	return &RoutingConfig{
		Version: 1,
		RoutingID: hash.New(id).String(),
		KeyPair:   kp,
	}, nil
}

// GenerateRouting generates a new routing structure
func GenerateRouting() (*RoutingConfig, error) {
	kt, err := bmcrypto.FindKeyType("ed25519")
	if err != nil {
		return nil, err
	}

	kp, err := internal.GenerateKeypairWithRandomSeed(kt)
	if err != nil {
		return nil, err
	}

	id := hex.EncodeToString(kp.PubKey.K.(ed25519.PublicKey))
	return &RoutingConfig{
		Version: 1,
		RoutingID: hash.New(id).String(),
		KeyPair:   kp,
	}, nil
}
