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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	pbkdfIterations = 100002
)

// DecryptContainer decrypts a container and fills the values in v.Store
func (v *Vault) DecryptContainer(container *vaultContainer) error {
	// Check if HMAC is correct
	hash := hmac.New(sha256.New, []byte(v.password))
	hash.Write(container.Data)
	if !bytes.Equal(hash.Sum(nil), container.Hmac) {
		return errIncorrectPassword
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key([]byte(v.password), container.Salt, pbkdfIterations, 32, sha256.New)
	aes256, err := aes.NewCipher(derivedAESKey)
	if err != nil {
		return err
	}

	// Decrypt vault data
	plainText := make([]byte, len(container.Data))
	ctr := cipher.NewCTR(aes256, container.Iv)
	ctr.XORKeyStream(plainText, container.Data)

	// store raw data. This makes editing through vault-edit tool easier
	v.RawData = plainText

	if container.Version == VersionV0 {
		// Unmarshal "old" style, with no organisations present
		var accounts []AccountInfo
		err = json.Unmarshal(plainText, &accounts)
		if err == nil {
			v.Store.Accounts = accounts
			v.Store.Organisations = []OrganisationInfo{}

			// Write back to disk in a newer format
			return v.Persist()
		}
	}

	// Version 1 has organisation info
	if container.Version == VersionV1 {
		err = json.Unmarshal(plainText, &v.Store)
		if err != nil {
			return err
		}
	}

	return nil
}

// EncryptContainer encrypts v.Store and returns the vault as encrypted JSON container
func (v *Vault) EncryptContainer() ([]byte, error) {
	// Generate 64 byte salt
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key([]byte(v.password), salt, pbkdfIterations, 32, sha256.New)
	aes256, err := aes.NewCipher(derivedAESKey)
	if err != nil {
		return nil, err
	}

	// Generate 32 byte IV
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	// Marshal and encrypt the data
	plainText, err := json.MarshalIndent(&v.Store, "", "  ")
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(plainText))
	ctr := cipher.NewCTR(aes256, iv)
	ctr.XORKeyStream(cipherText, plainText)

	// Generate HMAC based on the encrypted data (encrypt-then-mac?)
	hash := hmac.New(sha256.New, []byte(v.password))
	hash.Write(cipherText)

	// Generate the vault structure for disk
	return json.MarshalIndent(&vaultContainer{
		Version: VersionV1,
		Data:    cipherText,
		Salt:    salt,
		Iv:      iv,
		Hmac:    hash.Sum(nil),
	}, "", "  ")
}
