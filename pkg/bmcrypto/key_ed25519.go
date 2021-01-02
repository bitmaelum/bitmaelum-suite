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
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"errors"
	"io"
	"math/big"

	"github.com/vtolstov/jwt-go"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

var errCannotFetchScalar = errors.New("error while getting a random scalar")

// KeyEd25519 is the ed25519 keytype
type KeyEd25519 struct {
}

// NewEd25519Key will return the keytype of ed25519. There is only a single curve here.
func NewEd25519Key() KeyType {
	return &KeyEd25519{}
}

// CanEncrypt returns true if the key(type) is able to be used for encryption/decryption
func (k *KeyEd25519) CanEncrypt() bool {
	return false
}

// CanKeyExchange returns true if the key(type) is able to be used for key exchange
func (k *KeyEd25519) CanKeyExchange() bool {
	return true
}

// CanDualKeyExchange returns true if the key(type) is able to be used for a dual key exchange
func (k *KeyEd25519) CanDualKeyExchange() bool {
	return true
}

// String returns a string representation of the key type ("rsa", "ecdsa", "ed25519" etc)
func (k *KeyEd25519) String() string {
	return "ed25519"
}

// ParsePrivateKeyData will parse a string representation of a key and returns the given key
func (k *KeyEd25519) ParsePrivateKeyData(buf []byte) (interface{}, error) {
	return x509.ParsePKCS8PrivateKey(buf)
}

// ParsePrivateKeyInterface will parse a interface and returns the key representation
func (k *KeyEd25519) ParsePrivateKeyInterface(key interface{}) ([]byte, error) {
	switch key := key.(type) {
	case ed25519.PrivateKey:
		return x509.MarshalPKCS8PrivateKey(key)
	}

	return nil, errIncorrectKey
}

// GenerateKeyPair will generate a new keypair for this keytype. io.Reader can be deterministic if needed
func (k *KeyEd25519) GenerateKeyPair(r io.Reader) (*PrivKey, *PubKey, error) {
	// Read 192 bits
	randBuf := make([]byte, 24)
	_, err := io.ReadFull(r, randBuf)
	if err != nil {
		return nil, nil, err
	}

	// Stretch 192 bits to 256 bits
	rd := hkdf.New(sha256.New, randBuf[:], []byte{}, []byte{})
	expBuf := make([]byte, 32)
	_, err = io.ReadFull(rd, expBuf)
	if err != nil {
		return nil, nil, err
	}

	// Generate keypair
	pk := ed25519.NewKeyFromSeed(expBuf[:32])
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
func (k *KeyEd25519) JWTSignMethod() jwt.SigningMethod {
	return &SigningMethodEdDSA{}
}

// JWTHasValidSignMethod will return true when this keytype has been used for signing the token
func (k *KeyEd25519) JWTHasValidSignMethod(token *jwt.Token) bool {
	_, ok := token.Method.(*SigningMethodEdDSA)
	return ok
}

// Verify will verify the signature with the public key
func (k *KeyEd25519) Verify(key PubKey, message []byte, sig []byte) (bool, error) {
	return ed25519.Verify(key.K.(ed25519.PublicKey), message, sig), nil
}

// Sign will sign the given bytes with the private key
func (k *KeyEd25519) Sign(reader io.Reader, key PrivKey, message []byte) ([]byte, error) {
	return ed25519.Sign(key.K.(ed25519.PrivateKey), message), nil
}

// Encrypt will encrypt the given bytes with the public key. Will return the ciphertext, a transaction ID (if needed), the crypto used and an error
func (k *KeyEd25519) Encrypt(key PubKey, message []byte) ([]byte, string, string, error) {
	secret, txID, err := DualKeyExchange(key)
	if err != nil {
		return nil, "", "", err
	}

	encryptedMessage, err := MessageEncrypt(secret, message)

	return encryptedMessage, txID.ToHex(), "ed25519+aes", err
}

// Decrypt will decrypt the given bytes with the private key
func (k *KeyEd25519) Decrypt(key PrivKey, txID string, message []byte) ([]byte, error) {
	tx, err := TxIDFromString(txID)
	if err != nil {
		return nil, err
	}

	secret, ok, err := DualKeyGetSecret(key, *tx)
	if !ok || err != nil {
		return nil, err
	}

	return MessageDecrypt(secret, message)
}

// ParsePublicKeyData will parse a interface and returns the key representation
func (k *KeyEd25519) ParsePublicKeyData(buf []byte) (interface{}, error) {
	return x509.ParsePKIXPublicKey(buf)
}

// ParsePublicKeyInterface will parse a interface and returns the key representation
func (k *KeyEd25519) ParsePublicKeyInterface(key interface{}) ([]byte, error) {
	switch key := key.(type) {
	case ed25519.PublicKey:
		return x509.MarshalPKIXPublicKey(key)
	}

	return nil, errIncorrectKey
}

// KeyExchange allows for a key exchange (if possible in the keytype)
func (k *KeyEd25519) KeyExchange(privK PrivKey, pubK PubKey) ([]byte, error) {
	x25519priv := EdPrivToX25519(privK.K.(ed25519.PrivateKey))
	x25519pub := EdPubToX25519(pubK.K.(ed25519.PublicKey))

	return curve25519.X25519(x25519priv, x25519pub)
}

// EdPrivToX25519 converts a ed25519 PrivateKey to a X25519 Private Key
func EdPrivToX25519(privateKey ed25519.PrivateKey) []byte {
	h := sha512.New()
	_, _ = h.Write(privateKey[:32])
	digest := h.Sum(nil)
	h.Reset()

	/* From https://cr.yp.to/ecdh.html (I don't think this is really needed in this case)
	 * more info here: https://www.reddit.com/r/crypto/comments/66b3dp/how_do_is_a_curve25519_key_pair_generated/
	 */
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	return digest[:32]
}

var curve25519P, _ = new(big.Int).SetString("57896044618658097711785492504343953926634992332820282019728792003956564819949", 10)

// EdPubToX25519 converts a ed25519 Public Key to a X25519 Public Key
func EdPubToX25519(pk ed25519.PublicKey) []byte {
	// ed25519.PublicKey is a little endian representation of the y-coordinate,
	// with the most significant bit set based on the sign of the x-coordinate.
	bigEndianY := make([]byte, ed25519.PublicKeySize)
	for i, b := range pk {
		bigEndianY[ed25519.PublicKeySize-i-1] = b
	}
	bigEndianY[0] &= 127

	/* The Montgomery u-coordinate is derived through the bilinear map
	 *
	 *     u = (1 + y) / (1 - y)
	 *
	 * See https://blog.filippo.io/using-ed25519-keys-for-encryption.
	 */
	y := new(big.Int).SetBytes(bigEndianY)
	denom := big.NewInt(1)
	denom.ModInverse(denom.Sub(denom, y), curve25519P)
	u := y.Mul(y.Add(y, big.NewInt(1)), denom)
	u.Mod(u, curve25519P)

	out := make([]byte, curve25519.PointSize)
	uBytes := u.Bytes()
	for i, b := range uBytes {
		out[len(uBytes)-i-1] = b
	}

	return out
}

// DualKeyExchange allows for a ECIES key exchange
func (k *KeyEd25519) DualKeyExchange(pub PubKey) ([]byte, *TransactionID, error) {
	rs, err := generateRandomScalar()
	if err != nil {
		return nil, nil, errCannotFetchScalar
	}

	r := ed25519.NewKeyFromSeed(rs)
	R := r.Public()

	// Step 1: D = rA
	D, err := curve25519.X25519(
		EdPrivToX25519(r.Seed()),
		EdPubToX25519(pub.K.(ed25519.PublicKey)),
	)
	if err != nil {
		return nil, nil, err
	}

	f := hs(D)      // Step 2: f = Hs(D)
	P := f.Public() // Step 3-5: convert F into private Key (F=fG)

	return D, &TransactionID{
		P: P.(ed25519.PublicKey)[:32],
		R: R.(ed25519.PublicKey)[:32],
	}, nil
}
