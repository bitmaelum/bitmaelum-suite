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

package bmcrypto

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

var (
	errUnknownKeyType          = errors.New("unknown key type")
	errCannotUseForKeyExchange = errors.New("key type cannot be used for key exchange")
	errCannotUseForEncryption  = errors.New("this key type is not usable for encryption")
)

// Allows for easy mocking
var randReader io.Reader = rand.Reader

type ecdsaSignature struct {
	R, S *big.Int
}

// FindKeyType finds the keytype based on the given string
func FindKeyType(typ string) (KeyType, error) {
	for _, kt := range KeyTypes {
		if kt.String() == typ {
			return kt, nil
		}
	}

	return nil, errUnknownKeyType
}

// Sign a message based on the given key.
func Sign(key PrivKey, message []byte) ([]byte, error) {
	return key.Type.Sign(randReader, key, message)
}

// Verify if hash compares against the signature of the message
func Verify(key PubKey, message []byte, sig []byte) (bool, error) {
	return key.Type.Verify(key, message, sig)
}

// GenerateKeyPair generates a private/public keypair based on the given type
func GenerateKeyPair(kt KeyType) (*PrivKey, *PubKey, error) {
	if kt == nil {
		return nil, nil, errUnknownKeyType
	}

	return kt.GenerateKeyPair(randReader)
}

// KeyExchange exchange a message given the Private and other's Public Key
func KeyExchange(privK PrivKey, pubK PubKey) ([]byte, error) {
	if !privK.Type.CanKeyExchange() {
		return nil, errCannotUseForKeyExchange
	}

	return privK.Type.KeyExchange(privK, pubK)
}

// Encrypt a message with the given key
func Encrypt(pubKey PubKey, message []byte) ([]byte, string, string, error) {
	if !pubKey.Type.CanEncrypt() && !pubKey.Type.CanKeyExchange() {
		return nil, "", "", errCannotUseForEncryption
	}

	return pubKey.Type.Encrypt(pubKey, message)
}

// Decrypt a message with the given key
func Decrypt(key PrivKey, txID string, message []byte) ([]byte, error) {
	if !key.Type.CanEncrypt() && !key.Type.CanKeyExchange() {
		return nil, errCannotUseForEncryption
	}

	return key.Type.Decrypt(key, txID, message)
}
