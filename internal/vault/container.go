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

type EncryptedContainer struct {
	Type    string `json:"type"`
	Version int    `json:"version"`
	Data    []byte `json:"data"`
	Salt    []byte `json:"salt"`
	Iv      []byte `json:"iv"`
	Hmac    []byte `json:"hmac"`
}

// DecryptContainer decrypts an encrypted data container
func DecryptContainer(container *EncryptedContainer, pass string) ([]byte, error) {
	// Check if HMAC is correct
	hash := hmac.New(sha256.New, []byte(pass))
	hash.Write(container.Data)
	if !bytes.Equal(hash.Sum(nil), container.Hmac) {
		return nil, errIncorrectPassword
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key([]byte(pass), container.Salt, pbkdfIterations, 32, sha256.New)
	aes256, err := aes.NewCipher(derivedAESKey)
	if err != nil {
		return nil, err
	}

	// Decrypt vault data
	decryptedData := make([]byte, len(container.Data))
	ctr := cipher.NewCTR(aes256, container.Iv)
	ctr.XORKeyStream(decryptedData, container.Data)

	return decryptedData, nil
}

// EncryptContainer encrypts bytes into a data container
func EncryptContainer(data interface{}, pass, containerType string, version int) ([]byte, error) {
	// Generate 64 byte salt
	salt := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}

	// Generate key based on password
	derivedAESKey := pbkdf2.Key([]byte(pass), salt, pbkdfIterations, 32, sha256.New)
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
	plainText, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(plainText))
	ctr := cipher.NewCTR(aes256, iv)
	ctr.XORKeyStream(cipherText, plainText)

	// Generate HMAC based on the encrypted data (encrypt-then-mac?)
	hash := hmac.New(sha256.New, []byte(pass))
	hash.Write(cipherText)

	// Generate the encrypted container
	return json.MarshalIndent(&EncryptedContainer{
		Type:    containerType,
		Version: version,
		Data:    cipherText,
		Salt:    salt,
		Iv:      iv,
		Hmac:    hash.Sum(nil),
	}, "", "  ")
}
