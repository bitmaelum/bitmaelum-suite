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

package testing

import (
	"encoding/json"
	"io/ioutil"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// ReadKeyPair reads a path to a keypair
func ReadKeyPair(p string) (*bmcrypto.KeyPair, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	type jsonKeyType struct {
		PrivKey bmcrypto.PrivKey `json:"private_key"`
		PubKey  bmcrypto.PubKey  `json:"public_key"`
	}

	v := &jsonKeyType{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	return &bmcrypto.KeyPair{
		Generator:   "",
		FingerPrint: v.PubKey.Fingerprint(),
		PrivKey:     v.PrivKey,
		PubKey:      v.PubKey,
	}, nil
}

// ReadTestKey reads a path to a keypair and returns the keys
func ReadTestKey(p string) (*bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, nil, err
	}

	type jsonKeyType struct {
		PrivKey bmcrypto.PrivKey `json:"private_key"`
		PubKey  bmcrypto.PubKey  `json:"public_key"`
	}

	v := &jsonKeyType{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, nil, err
	}

	return &v.PrivKey, &v.PubKey, nil
}

// ReadTestFile reads a file
func ReadTestFile(p string) []byte {
	data, _ := ioutil.ReadFile(p)
	return data
}

// ReadJSON reads a json file and returns it in the given interface
func ReadJSON(p string, v interface{}) error {
	data, _ := ioutil.ReadFile(p)

	return json.Unmarshal(data, v)
}
