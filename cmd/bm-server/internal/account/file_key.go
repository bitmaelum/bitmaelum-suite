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

package account

import (
	"encoding/json"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
)

// Store the public key for this account
func (r *fileRepo) StoreKey(addr hash.Hash, key bmcrypto.PubKey) error {
	// Lock our key file for writing
	lockfilePath := r.getPath(addr, keysFile+".lock")
	lock, err := r.lockFile(lockfilePath)
	if err != nil {
		return err
	}

	err = r.tryLockFile(lock)
	if err != nil {
		return err
	}

	defer func() {
		_ = r.unLockFile(lock)
	}()

	// Read keys
	pk := &PubKeys{}
	err = r.fetchJSON(addr, keysFile, pk)
	if err != nil && os.IsNotExist(err) {
		err = r.createPubKeyFile(addr)
	}
	if err != nil {
		return err
	}

	// Add new key
	pk.PubKeys = append(pk.PubKeys, key)

	// Convert back to string
	data, err := json.MarshalIndent(pk, "", "  ")
	if err != nil {
		return err
	}

	// And store
	return r.store(addr, keysFile, data)
}

// Retrieve the public keys for this account
func (r *fileRepo) FetchKeys(addr hash.Hash) ([]bmcrypto.PubKey, error) {
	pk := &PubKeys{}
	err := r.fetchJSON(addr, keysFile, pk)
	if err != nil && os.IsNotExist(err) {
		err = r.createPubKeyFile(addr)
	}
	if err != nil {
		return nil, err
	}

	return pk.PubKeys, nil
}

// Create file because it doesn't exist yet
func (r *fileRepo) createPubKeyFile(addr hash.Hash) error {
	p := r.getPath(addr, keysFile)
	f, err := r.fs.Create(p)
	if err != nil {
		logrus.Errorf("Error while creating file %s: %s", p, err)
		return err
	}
	return f.Close()
}
