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

package bmcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
)

// keyGenerator can be used for mocking purposes
var (
	keyGenerator = randomKeyGenerator
)

const (
	keySize = 32
)

// JSONEncrypt encrypts a structure that is marshalled to JSON
func JSONEncrypt(key []byte, data interface{}) ([]byte, error) {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return MessageEncrypt(key, plaintext)
}

// JSONDecrypt decrypts data back from a encrypted marshalled JSON structure
func JSONDecrypt(key []byte, ciphertext []byte, v interface{}) error {
	plaintext, err := MessageDecrypt(key, ciphertext)
	if err != nil {
		return err
	}

	return json.Unmarshal(plaintext, &v)
}

// MessageEncrypt encrypts a binary message
func MessageEncrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := keyGenerator(aead.NonceSize())
	if err != nil {
		return nil, err
	}

	return append(nonce, aead.Seal(nil, nonce, plaintext, nil)...), nil
}

// MessageDecrypt decrypts a binary message
func MessageDecrypt(key []byte, message []byte) ([]byte, error) {
	// Key should be 32byte (256bit)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()
	if len(message) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := message[:nonceSize], message[nonceSize:]
	return aead.Open(nil, nonce, ciphertext, nil)
}

// CreateCatalogKey generates a new key that is used for encrypting a catalog
func CreateCatalogKey() ([]byte, error) {
	return keyGenerator(keySize)
}

// CatalogEncrypt encrypts a catalog with a random key.
func CatalogEncrypt(catalog interface{}) ([]byte, []byte, error) {
	catalogKey, err := keyGenerator(keySize)
	if err != nil {
		return nil, nil, err
	}

	ciphertext, err := JSONEncrypt(catalogKey, catalog)
	if err != nil {
		return nil, nil, err
	}

	return catalogKey, ciphertext, nil
}

// CatalogDecrypt decrypts a catalog with the given key
func CatalogDecrypt(key, data []byte, catalog interface{}) error {
	err := JSONDecrypt(key, data, catalog)
	if err != nil {
		return err
	}

	return nil
}

// randomKeyGenrator generates keys for catalog encryption.
func randomKeyGenerator(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, key)

	return key, err
}

// GenerateIvAndKey generate a random IV and key
func GenerateIvAndKey() ([]byte, []byte, error) {
	iv := make([]byte, 16)
	n, err := rand.Read(iv)
	if n != 16 || err != nil {
		return nil, nil, err
	}

	key := make([]byte, 32)
	n, err = rand.Read(key)
	if n != 32 || err != nil {
		return nil, nil, err
	}

	return iv, key, nil
}

// GetAesEncryptorReader returns a reader that automatically encrypts reader blocks through CFB stream
func GetAesEncryptorReader(iv []byte, key []byte, r io.Reader) (io.Reader, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("IV is not of blocksize")
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, err
}

// GetAesDecryptorReader returns a reader that automatically decrypts reader blocks through CFB stream
func GetAesDecryptorReader(iv []byte, key []byte, r io.Reader) (io.Reader, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	if len(iv) != block.BlockSize() {
		return nil, errors.New("IV is not of blocksize")
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, err
}
