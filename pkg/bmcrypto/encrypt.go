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
	"crypto/rsa"
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
)

// Encrypt a message with the given key
func Encrypt(pubKey PubKey, message []byte) ([]byte, string, string, error) {
	if !pubKey.CanEncrypt() && !pubKey.CanKeyExchange() {
		return nil, "", "", errors.New("this key type is not usable for encryption")
	}

	switch pubKey.Type {
	case KeyTypeRSA:
		encryptedMessage, err := encryptRsa(pubKey.K.(*rsa.PublicKey), message)
		return encryptedMessage, "", "rsa+aes256gcm", err
		// TODO: Implement KeyTypeECDSA
	case KeyTypeED25519:
		return encryptED25519(pubKey, message)
	}

	return nil, "", "", errors.New("encryption not implemented for" + pubKey.Type)
}

// Decrypt a message with the given key
func Decrypt(key PrivKey, txID string, message []byte) ([]byte, error) {
	if !key.CanEncrypt() && !key.CanKeyExchange() {
		return nil, errors.New("this key type is not usable for encryption")
	}

	switch key.Type {
	case KeyTypeRSA:
		return decryptRsa(key.K.(*rsa.PrivateKey), message)
	// TODO: Implement KeyTypeECDSA
	case KeyTypeED25519:
		return decryptED25519(key, txID, message)

	}

	return nil, errors.New("encryption not implemented for" + key.Type)
}

func encryptRsa(key *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, key, message)
}

func decryptRsa(key *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, message)
}

func encryptED25519(pubKey PubKey, message []byte) ([]byte, string, string, error) {
	secret, txID, err := DualKeyExchange(pubKey)
	if err != nil {
		return nil, "", "", err
	}

	encryptedMessage, err := encrypt.MessageEncrypt(secret, message)
	return encryptedMessage, txID.ToHex(), "ed25519+aes256gcm", err
}

func decryptED25519(privKey PrivKey, txIDString string, message []byte) ([]byte, error) {
	txID, err := TxIDFromString(txIDString)
	if err != nil {
		return nil, err
	}

	secret, ok, err := DualKeyGetSecret(privKey, *txID)
	if !ok || err != nil {
		return nil, err
	}

	return encrypt.MessageDecrypt(secret, message)
}
