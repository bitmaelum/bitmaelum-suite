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

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

var errMigrationError = errors.New("error while migrating vault to the latest version")
var errGenerateStoreKey = errors.New("error while generating store keys")

// DecryptContainer decrypts a container and fills the values in v.Store
func (v *Vault) DecryptContainer(container *EncryptedContainer) error {

	var err error
	v.RawData, err = DecryptContainer(container, v.password)
	if err != nil {
		return err
	}

	// Migrate vault to latest version if needed
	store, err := MigrateVault(v.RawData, container.Version)
	if err != nil {
		return errMigrationError
	}
	v.Store = *store

	// Save vault if it wasn't the latest version already
	if container.Version != LatestVaultVersion {
		err = v.Persist()
		if err != nil {
			return errMigrationError
		}
	}

	// Check if we have store keys
	updated := false
	for i := range v.Store.Accounts {
		if v.Store.Accounts[i].StoreKey == nil {
			kt, err := bmcrypto.FindKeyType("ed25519")
			if err != nil {
				return errGenerateStoreKey
			}

			kp, err := bmcrypto.GenerateKeypairWithRandomSeed(kt)
			if err != nil {
				return errGenerateStoreKey
			}
			v.Store.Accounts[i].StoreKey = kp
			updated = true
		}
	}

	// Save vault if store keys are created
	if updated {
		err = v.Persist()
		if err != nil {
			return errGenerateStoreKey
		}
	}

	return nil
}

// EncryptContainer encrypts v.Store and returns the vault as encrypted JSON container
func (v *Vault) EncryptContainer() ([]byte, error) {
	return EncryptContainer(&v.Store, v.password, "vault", LatestVaultVersion)
}
