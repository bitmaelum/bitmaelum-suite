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

package bmcrypto

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"io"

	deterministicRsaKeygen "github.com/cloudflare/gokey/rsa"
	"github.com/vtolstov/jwt-go"
	"golang.org/x/crypto/hkdf"
)

// Taken from https://github.com/cloudflare/gokey/blob/6bb7290160583cf1fd7cdcb5726093a00dd23c25/csprng.go#L18
type devZero struct{}

func (dz devZero) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// Hash used for OAEP encryption/decryption
var oaepHash = sha256.New()

// KeyRsa is the RSA keytype
type KeyRsa struct {
	Type    string
	BitSize int
}

// NewRsaKey will return a keytype based on the number of bits
func NewRsaKey(bits int) KeyType {
	k := &KeyRsa{}

	k.BitSize = bits
	k.Type = fmt.Sprintf("rsa%d", bits)

	if bits == 2048 {
		// rsa2048 is actually plain "rsa"
		k.Type = "rsa"
	}

	return k
}

// CanEncrypt returns true if the key(type) is able to be used for encryption/decryption
func (k *KeyRsa) CanEncrypt() bool {
	return true
}

// CanKeyExchange returns true if the key(type) is able to be used for key exchange
func (k *KeyRsa) CanKeyExchange() bool {
	return false
}

// CanDualKeyExchange returns true if the key(type) is able to be used for a dual key exchange
func (k *KeyRsa) CanDualKeyExchange() bool {
	return false
}

// String returns a string representation of the key type ("rsa", "ecdsa", "ed25519" etc)
func (k *KeyRsa) String() string {
	return k.Type
}

// ParsePrivateKeyData will parse a string representation of a key and returns the given key
func (k *KeyRsa) ParsePrivateKeyData(buf []byte) (interface{}, error) {
	pk, err := x509.ParsePKCS1PrivateKey(buf)
	if err == nil {
		return pk, nil
	}

	return x509.ParsePKCS8PrivateKey(buf)
}

// ParsePrivateKeyInterface will parse a interface and returns the key representation
func (k *KeyRsa) ParsePrivateKeyInterface(key interface{}) ([]byte, error) {
	switch key := key.(type) {
	case *rsa.PrivateKey:
		if key.Size()*8 == k.BitSize {
			return x509.MarshalPKCS8PrivateKey(key)
		}
	}

	return nil, errors.New("incorrect key")
}

// GenerateKeyPair will generate a new keypair for this keytype. io.Reader can be deterministic if needed
func (k *KeyRsa) GenerateKeyPair(r io.Reader) (*PrivKey, *PubKey, error) {
	// @TODO: we should check this:
	// We have a reader that outputs 24bytes (192 bits), and stretches this to 256 bits.
	// Then, create a new cipher stream, but how long is this stream?? Is it enough to create
	// deterministic RSA keys?

	// Read 192 bits
	randBuf := make([]byte, 24)
	_, err := io.ReadFull(r, randBuf)
	if err != nil {
		return nil, nil, err
	}

	// Stretch 192 bits to 256 bits
	rd := hkdf.New(sha256.New, randBuf, []byte{}, []byte{})
	expBuf := make([]byte, 32)
	_, err = io.ReadFull(rd, expBuf)
	if err != nil {
		return nil, nil, err
	}

	// Generate keypair
	// Taken from https://github.com/cloudflare/gokey/blob/6bb7290160583cf1fd7cdcb5726093a00dd23c25/csprng.go#L56
	block, _ := aes.NewCipher(expBuf[:32])
	stream := cipher.NewCTR(block, make([]byte, 16))

	randReader := cipher.StreamReader{S: stream, R: devZero{}}
	privRSAKey, err := deterministicRsaKeygen.GenerateKey(randReader, k.BitSize)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := PrivateKeyFromInterface(k, privRSAKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := PublicKeyFromInterface(k, privRSAKey.Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

// JWTSignMethod will return the signing method used for this keytype
func (k *KeyRsa) JWTSignMethod() jwt.SigningMethod {
	return jwt.SigningMethodRS256
}

// JWTHasValidSignMethod will return true when this keytype has been used for signing the token
func (k *KeyRsa) JWTHasValidSignMethod(token *jwt.Token) bool {
	_, ok := token.Method.(*jwt.SigningMethodRSA)
	return ok
}

// Verify will verify the signature with the public key
func (k *KeyRsa) Verify(key PubKey, message []byte, sig []byte) (bool, error) {
	h := crypto.SHA256.New()
	h.Write(message)
	hash := h.Sum(nil)

	err := rsa.VerifyPKCS1v15(key.K.(*rsa.PublicKey), crypto.SHA256, hash[:], sig)
	return err == nil, err
}

// Sign will sign the given bytes with the private key
func (k *KeyRsa) Sign(_ io.Reader, key PrivKey, message []byte) ([]byte, error) {
	h := crypto.SHA256.New()
	h.Write(message)
	hash := h.Sum(nil)

	return rsa.SignPKCS1v15(randReader, key.K.(*rsa.PrivateKey), crypto.SHA256, hash[:])
}

// Encrypt will encrypt the given message with the public key
func (k *KeyRsa) Encrypt(pubKey PubKey, msg []byte) ([]byte, *EncryptionSettings, error) {
	data, err := rsa.EncryptOAEP(oaepHash, rand.Reader, pubKey.K.(*rsa.PublicKey), msg, nil)
	if err != nil {
		return nil, nil, err
	}

	return data, &EncryptionSettings{
		Type:          RsaOAEP,
		TransactionID: "",
	}, nil
}

// Decrypt will decrypt the given ciphertext with the private key
func (k *KeyRsa) Decrypt(key PrivKey, settings *EncryptionSettings, ciphertext []byte) ([]byte, error) {
	switch settings.Type {
	default:
		fallthrough
	case Rsav15:
		return rsa.DecryptPKCS1v15(rand.Reader, key.K.(*rsa.PrivateKey), ciphertext)
	case RsaOAEP:
		return rsa.DecryptOAEP(oaepHash, rand.Reader, key.K.(*rsa.PrivateKey), ciphertext, nil)
	}
}

// ParsePublicKeyInterface will parse a interface and returns the key representation
func (k *KeyRsa) ParsePublicKeyInterface(key interface{}) ([]byte, error) {
	switch key := key.(type) {
	case *rsa.PublicKey:
		if key.Size()*8 == k.BitSize {
			return x509.MarshalPKIXPublicKey(key)
		}
	}

	return nil, errIncorrectKey
}

// ParsePublicKeyData will parse a interface and returns the key representation
func (k *KeyRsa) ParsePublicKeyData(buf []byte) (interface{}, error) {
	return x509.ParsePKIXPublicKey(buf)
}

// KeyExchange allows for a key exchange (if possible in the keytype)
func (k *KeyRsa) KeyExchange(_ PrivKey, _ PubKey) ([]byte, error) {
	return nil, errors.New("cannot exchange with RSA")
}

// DualKeyExchange allows for a ECIES key exchange
func (k *KeyRsa) DualKeyExchange(_ PubKey) ([]byte, *TransactionID, error) {
	return nil, nil, errors.New("cannot dual exchange with RSA")
}
