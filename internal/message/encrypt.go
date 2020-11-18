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

package message

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// Encrypt a message with the given key
func Encrypt(pubKey bmcrypto.PubKey, message []byte) ([]byte, string, string, error) {
	if !pubKey.Type.CanEncrypt() && !pubKey.Type.CanKeyExchange() {
		return nil, "", "", errors.New("this key type is not usable for encryption")
	}

	return pubKey.Type.Encrypt(pubKey, message)
}

// Decrypt a message with the given key
func Decrypt(key bmcrypto.PrivKey, txID string, message []byte) ([]byte, error) {
	if !key.Type.CanEncrypt() && !key.Type.CanKeyExchange() {
		return nil, errors.New("this key type is not usable for encryption")
	}

	return key.Type.Decrypt(key, txID, message)
}
