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
	"encoding/json"
)

// DecryptContainer decrypts a container and fills the values in v.Store
func (v *Vault) DecryptContainer(container *EncryptedContainer) error {

	var err error
	v.RawData, err = DecryptContainer(container, v.password)
	if err != nil {
		return err
	}

	if container.Version == VersionV0 {
		// Unmarshal "old" style, with no organisations present
		var accounts []AccountInfo
		err = json.Unmarshal(v.RawData, &accounts)
		if err == nil {
			v.Store.Accounts = accounts
			v.Store.Organisations = []OrganisationInfo{}

			// Write back to disk in a newer format
			return v.Persist()
		}
	}

	// Version 1 has organisation info
	if container.Version == VersionV1 {
		err = json.Unmarshal(v.RawData, &v.Store)
		if err != nil {
			return err
		}
	}

	return nil
}

// EncryptContainer encrypts v.Store and returns the vault as encrypted JSON container
func (v *Vault) EncryptContainer() ([]byte, error) {
	return EncryptContainer(&v.Store, v.password, "vault", VersionV1)
}
