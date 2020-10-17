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
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/sha256"
	"errors"
	"io"
	"strings"

	deterministicRsaKeygen "github.com/cloudflare/gokey/rsa"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/hkdf"
)

//Taken from https://github.com/cloudflare/gokey/blob/6bb7290160583cf1fd7cdcb5726093a00dd23c25/csprng.go#L18
type devZero struct{}

func (dz devZero) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// GenerateRSAKeypairFromMnemonic generates a keypair based on the given mnemonic
func GenerateRSAKeypairFromMnemonic(mnemonic string) (*PrivKey, *PubKey, error) {
	e, err := bip39.MnemonicToByteArray(mnemonic, true)
	if err != nil {
		return nil, nil, err
	}

	return genRSAKey(e)
}

// GenerateRSAKeypairWithMnemonic generates a mnemonic, and a RSA keypair that can be generated through the same mnemonic again.
func GenerateRSAKeypairWithMnemonic() (string, *PrivKey, *PubKey, error) {
	// Generate large enough random string
	e, err := bip39.NewEntropy(192)
	if err != nil {
		return "", nil, nil, err
	}

	// Generate Mnemonic words
	mnemonic, err := bip39.NewMnemonic(e)
	if err != nil {
		return "", nil, nil, err
	}

	privKey, pubKey, err := genRSAKey(e)
	if err != nil {
		return "", nil, nil, err
	}

	return "rsa " + mnemonic, privKey, pubKey, nil
}

// GenerateKeypairFromMnemonic generates a keypair based on the given mnemonic
func GenerateKeypairFromMnemonic(mnemonic string) (*PrivKey, *PubKey, error) {
	// Check if it's a RSA mnemonic
	words := strings.SplitN(mnemonic, " ", 2)
	if strings.ToLower(words[0]) == "rsa" {
		return GenerateRSAKeypairFromMnemonic(strings.Join(words[1:], " "))
	}

	e, err := bip39.MnemonicToByteArray(mnemonic, true)
	if err != nil {
		return nil, nil, err
	}

	return genKey(e)
}

// GenerateKeypairWithMnemonic generates a mnemonic, and a keypair that can be generated through the same mnemonic again.
func GenerateKeypairWithMnemonic(kt string) (string, *PrivKey, *PubKey, error) {
	switch kt {
	case KeyTypeRSA:
		return GenerateRSAKeypairWithMnemonic()
	case KeyTypeED25519:
		return GenerateED25519KeypairWithMnemonic()
	default:
		return "", nil, nil, errors.New("key type not supported")
	}
}

// GenerateED25519KeypairWithMnemonic generates a mnemonic, and a ED25519 keypair that can be generated through the same mnemonic again.
func GenerateED25519KeypairWithMnemonic() (string, *PrivKey, *PubKey, error) {
	// Generate large enough random string
	e, err := bip39.NewEntropy(192)
	if err != nil {
		return "", nil, nil, err
	}

	// Generate Mnemonic words
	mnemonic, err := bip39.NewMnemonic(e)
	if err != nil {
		return "", nil, nil, err
	}

	privKey, pubKey, err := genKey(e)
	if err != nil {
		return "", nil, nil, err
	}

	return mnemonic, privKey, pubKey, nil
}

func genRSAKey(e []byte) (*PrivKey, *PubKey, error) {
	// Stretch 192 bits to 256 bits
	rd := hkdf.New(sha256.New, e, []byte{}, []byte{})
	expbuf := make([]byte, 32)
	_, err := io.ReadFull(rd, expbuf)
	if err != nil {
		return nil, nil, err
	}

	// Generate keypair
	// Taken from https://github.com/cloudflare/gokey/blob/6bb7290160583cf1fd7cdcb5726093a00dd23c25/csprng.go#L56
	block, _ := aes.NewCipher(expbuf[:32])
	stream := cipher.NewCTR(block, make([]byte, 16))

	randReader := cipher.StreamReader{S: stream, R: devZero{}}
	privRSAKey, err := deterministicRsaKeygen.GenerateKey(randReader, rsaBits)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := NewPrivKeyFromInterface(privRSAKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := NewPubKeyFromInterface(privRSAKey.Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

func genKey(e []byte) (*PrivKey, *PubKey, error) {
	// Stretch 192 bits to 256 bits
	rd := hkdf.New(sha256.New, e, []byte{}, []byte{})
	expbuf := make([]byte, 32)
	_, err := io.ReadFull(rd, expbuf)
	if err != nil {
		return nil, nil, err
	}

	// Generate keypair
	r := ed25519.NewKeyFromSeed(expbuf[:32])
	privKey, err := NewPrivKeyFromInterface(r)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := NewPubKeyFromInterface(r.Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}
