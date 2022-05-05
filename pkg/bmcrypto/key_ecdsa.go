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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"io"

	"github.com/vtolstov/jwt-go"
)

// KeyEcdsa is a keytype for elliptic curve
type KeyEcdsa struct {
	Curve elliptic.Curve
}

// NewEcdsaKey creates a new keytype based on the given curve
func NewEcdsaKey(curve elliptic.Curve) KeyType {
	return &KeyEcdsa{
		Curve: curve,
	}
}

// CanEncrypt returns true if the key(type) is able to be used for encryption/decryption
func (k *KeyEcdsa) CanEncrypt() bool {
	return false
}

// CanKeyExchange returns true if the key(type) is able to be used for key exchange
func (k *KeyEcdsa) CanKeyExchange() bool {
	return true
}

// CanDualKeyExchange returns true if the key(type) is able to be used for a dual key exchange
func (k *KeyEcdsa) CanDualKeyExchange() bool {
	return false
}

// String returns a string representation of the key type ("rsa", "ecdsa", "ed25519" etc)
func (k *KeyEcdsa) String() string {
	return "ecdsa"
}

// ParsePrivateKeyData will parse a string representation of a key and returns the given key
func (k *KeyEcdsa) ParsePrivateKeyData(buf []byte) (interface{}, error) {
	return x509.ParsePKCS8PrivateKey(buf)
}

// ParsePrivateKeyInterface will parse a interface and returns the key representation
func (k *KeyEcdsa) ParsePrivateKeyInterface(key interface{}) ([]byte, error) {
	switch key := key.(type) {
	case *ecdsa.PrivateKey:
		return x509.MarshalPKCS8PrivateKey(key)
	}

	return nil, errIncorrectKey
}

// GenerateKeyPair will generate a new keypair for this keytype. io.Reader can be deterministic if needed
func (k *KeyEcdsa) GenerateKeyPair(r io.Reader) (*PrivKey, *PubKey, error) {
	pk, err := ecdsa.GenerateKey(k.Curve, r)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := PrivateKeyFromInterface(k, pk)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := PublicKeyFromInterface(k, pk.Public())
	if err != nil {
		return nil, nil, err
	}
	return privKey, pubKey, nil
}

// JWTSignMethod will return the signing method used for this keytype
func (k *KeyEcdsa) JWTSignMethod() jwt.SigningMethod {
	return jwt.SigningMethodES384
}

// JWTHasValidSignMethod will return true when this keytype has been used for signing the token
func (k *KeyEcdsa) JWTHasValidSignMethod(token *jwt.Token) bool {
	_, ok := token.Method.(*jwt.SigningMethodECDSA)
	return ok
}

// Verify will verify the signature with the public key
func (k *KeyEcdsa) Verify(key PubKey, message []byte, sig []byte) (bool, error) {
	ecdsaSig := ecdsaSignature{}
	_, err := asn1.Unmarshal(sig, &ecdsaSig)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(key.K.(*ecdsa.PublicKey), message, ecdsaSig.R, ecdsaSig.S), nil
}

// Sign will sign the given bytes with the private key
func (k *KeyEcdsa) Sign(_ io.Reader, key PrivKey, message []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(randReader, key.K.(*ecdsa.PrivateKey), message)
	if err != nil {
		return nil, err
	}

	sig := ecdsaSignature{
		R: r,
		S: s,
	}

	return asn1.Marshal(sig)
}

// Encrypt will encrypt the given msg with the public key.
func (k *KeyEcdsa) Encrypt(key PubKey, msg []byte) ([]byte, *EncryptionSettings, error) {
	secret, txID, err := DualKeyExchange(key)
	if err != nil {
		return nil, nil, err
	}

	encryptedMessage, err := MessageEncrypt(secret, msg)
	if err != nil {
		return nil, nil, err
	}

	return encryptedMessage, &EncryptionSettings{
		Type:          EcdsaAES,
		TransactionID: txID.ToHex(),
	}, nil
}

// Decrypt will decrypt the given bytes with the private key
func (k *KeyEcdsa) Decrypt(key PrivKey, settings *EncryptionSettings, cipherText []byte) ([]byte, error) {
	if settings.Type != EcdsaAES {
		return nil, errors.New("cannot decrypt this encryption type")
	}

	tx, err := TxIDFromString(settings.TransactionID)
	if err != nil {
		return nil, err
	}

	secret, ok, err := DualKeyGetSecret(key, *tx)
	if !ok || err != nil {
		return nil, err
	}

	return MessageDecrypt(secret, cipherText)
}

// ParsePublicKeyData will parse a interface and returns the key representation
func (k *KeyEcdsa) ParsePublicKeyData(buf []byte) (interface{}, error) {
	return x509.ParsePKIXPublicKey(buf)
}

// ParsePublicKeyInterface will parse a interface and returns the key representation
func (k *KeyEcdsa) ParsePublicKeyInterface(key interface{}) ([]byte, error) {
	switch key := key.(type) {
	case *ecdsa.PublicKey:
		return x509.MarshalPKIXPublicKey(key)
	}

	return nil, errIncorrectKey
}

// KeyExchange allows for a key exchange (if possible in the keytype)
func (k *KeyEcdsa) KeyExchange(privK PrivKey, pubK PubKey) ([]byte, error) {
	ke, _ := pubK.K.(*ecdsa.PublicKey).Curve.ScalarMult(
		pubK.K.(*ecdsa.PublicKey).X,
		pubK.K.(*ecdsa.PublicKey).Y,
		privK.K.(*ecdsa.PrivateKey).D.Bytes(),
	)

	// Get the length of the key
	keyLen := k.Curve.Params().BitSize / 8

	b := ke.Bytes()
	if len(b) == keyLen {
		// Length is 48 bytes, so we can return as-is
		return b, nil
	}

	// Sanity check
	if keyLen > len(b) {
		panic("key length is larger. This would mean we might have ended up with a truncated key.")
	}

	// Make sure we zero-extend the result (big.Int) to 32 bytes (big endian)
	var ret = make([]byte, keyLen)
	copy(ret[keyLen-len(b):], b)
	return ret[:], nil
}

// DualKeyExchange allows for a ECIES key exchange
func (k *KeyEcdsa) DualKeyExchange(pub PubKey) ([]byte, *TransactionID, error) {
	return nil, nil, errCannotuseForDualKeyExchange
}
