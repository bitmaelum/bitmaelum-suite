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
	"bytes"
	"crypto/rand"
	"encoding/hex"
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

// KeyPair is a structure with key information
type KeyPair struct {
	Generator   string  `json:"generator"`   // The generator string that will generate the given keypair
	FingerPrint string  `json:"fingerprint"` // The sha1 fingerprint for this key
	PrivKey     PrivKey `json:"priv_key"`    // PEM encoded private key
	PubKey      PubKey  `json:"pub_key"`     // PEM encoded public key
}

// CryptoType is a type that defines which crypto is used
type CryptoType string

// Different crypto types used
const (
	Rsav15     CryptoType = "rsa"         // PKCS1v1.5 (obsolete)
	RsaOAEP    CryptoType = "rsa/oaep"    // PKCS1 OAEP (default)
	Ed25519AES CryptoType = "ed25519+aes" // ED25519 ECEIS + AES
	EcdsaAES   CryptoType = "ecdsa+aes"   // ECDSA ECEIS + AES
)

// EncryptionSettings is a structure that passes information about the given encryption
type EncryptionSettings struct {
	Type          CryptoType
	TransactionID string
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
func Encrypt(pubKey PubKey, msg []byte) ([]byte, *EncryptionSettings, error) {
	if !pubKey.Type.CanEncrypt() && !pubKey.Type.CanKeyExchange() {
		return nil, nil, errCannotUseForEncryption
	}

	return pubKey.Type.Encrypt(pubKey, msg)
}

// Decrypt a message with the given key
func Decrypt(key PrivKey, settings *EncryptionSettings, ciphertext []byte) ([]byte, error) {
	if !key.Type.CanEncrypt() && !key.Type.CanKeyExchange() {
		return nil, errCannotUseForEncryption
	}

	return key.Type.Decrypt(key, settings, ciphertext)
}

// CreateKeypair create a new keypair
func CreateKeypair(kt KeyType, seed []byte) (*KeyPair, error) {
	privKey, pubKey, err := kt.GenerateKeyPair(bytes.NewReader(seed))
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		Generator:   hex.EncodeToString(seed),
		FingerPrint: pubKey.Fingerprint(),
		PrivKey:     *privKey,
		PubKey:      *pubKey,
	}, nil
}
