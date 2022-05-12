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
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vtolstov/jwt-go"
)

// KeyType is an interface that each key type should implement.
type KeyType interface {
	// CanEncrypt returns true if the key(type) is able to be used for encryption/decryption
	CanEncrypt() bool
	// CanKeyExchange returns true if the key(type) is able to be used for key exchange
	CanKeyExchange() bool
	// CanDualKeyExchange returns true if the key(type) is able to be used for a dual key exchange
	CanDualKeyExchange() bool

	// String returns a string representation of the key type ("rsa", "ecdsa", "ed25519" etc)
	String() string

	// ParsePrivateKeyData will parse a string representation of a key and returns the given key
	ParsePrivateKeyData([]byte) (interface{}, error)
	// ParsePrivateKeyInterface will parse a interface and returns the key representation
	ParsePrivateKeyInterface(interface{}) ([]byte, error)
	// ParsePublicKeyData will parse a interface and returns the key representation
	ParsePublicKeyData([]byte) (interface{}, error)
	// ParsePublicKeyInterface will parse a interface and returns the key representation
	ParsePublicKeyInterface(interface{}) ([]byte, error)

	// GenerateKeyPair will generate a new keypair for this keytype. io.Reader can be deterministic if needed
	GenerateKeyPair(io.Reader) (*PrivKey, *PubKey, error)

	// JWTSignMethod will return the signing method used for this keytype
	JWTSignMethod() jwt.SigningMethod
	// JWTHasValidSignMethod will return true when this keytype has been used for signing the token
	JWTHasValidSignMethod(*jwt.Token) bool

	// Encrypt will encrypt the given message with the public key.
	Encrypt(PubKey, []byte) ([]byte, *EncryptionSettings, error)
	// Decrypt will decrypt the given ciphertext with the private key
	Decrypt(PrivKey, *EncryptionSettings, []byte) ([]byte, error)

	// Sign will sign the given bytes with the private key
	Sign(io.Reader, PrivKey, []byte) ([]byte, error)
	// Verify will verify the signature with the public key
	Verify(PubKey, []byte, []byte) (bool, error)

	// KeyExchange allows for a key exchange (if possible in the keytype)
	KeyExchange(privK PrivKey, pubK PubKey) ([]byte, error)
	// DualKeyExchange allows for a ECIES key exchange
	DualKeyExchange(_ PubKey) ([]byte, *TransactionID, error)
}

var (
	errIncorrectKeyFormat = errors.New("incorrect key format")
	errUnsupportedKeyType = errors.New("unsupported key type")
	errIncorrectKey       = errors.New("incorrec key")
)

// KeyTypes is a list of all keytypes available
var KeyTypes = []KeyType{
	NewRsaKey(2048),              // RSA 2048 bits
	NewRsaKey(4096),              // RSA 4096 bits
	NewEd25519Key(),              // ED25519
	NewEcdsaKey(elliptic.P384()), // Ecdsa (P384)
}

// PrivKey is a structure containing a private key in multiple formats
type PrivKey struct {
	Type KeyType     // structure of the key
	S    string      // String representation <type> <PEM key>
	B    []byte      // Byte representation of string
	K    interface{} // Key interface{}
}

// Strings returns the key in a textual representation
func (pk *PrivKey) String() string {
	return fmt.Sprintf("%s %s", pk.Type.String(), pk.S)
}

// MarshalJSON marshals a key into bytes
func (pk *PrivKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pk.String())
}

// UnmarshalJSON unmarshals bytes into a key
func (pk *PrivKey) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	pkCopy, err := PrivateKeyFromString(s)
	if err != nil {
		return err
	}

	*pk = *pkCopy
	return err
}

// PrivateKeyFromString creates a new private key based on the string data "<type> <key>"
func PrivateKeyFromString(data string) (*PrivKey, error) {
	if !strings.Contains(data, " ") {
		return nil, errIncorrectKeyFormat
	}

	// <type> <data>
	parts := strings.SplitN(data, " ", 2)
	if len(parts) != 2 {
		return nil, errIncorrectKeyFormat
	}

	pk := &PrivKey{}

	// Find the correct key type
	var err error
	pk.Type, err = FindKeyType(parts[0])
	if err != nil {
		return nil, errUnsupportedKeyType
	}

	// Set values
	pk.S = strings.TrimSpace(parts[1])
	pk.B = []byte(pk.S)

	// Decode base64 before we parse to key
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(pk.B)))
	n, err := base64.StdEncoding.Decode(buf, pk.B)
	if err != nil {
		return nil, errIncorrectKeyFormat
	}

	// Parse key
	pk.K, err = pk.Type.ParsePrivateKeyData(buf[:n])

	if err != nil {
		return nil, errIncorrectKeyFormat
	}

	return pk, nil
}

// PrivateKeyFromInterface creates a new key based on an interface{} (like rsa.PrivateKey)
func PrivateKeyFromInterface(kt KeyType, key interface{}) (*PrivKey, error) {
	buf, err := kt.ParsePrivateKeyInterface(key)
	if err != nil {
		return nil, errIncorrectKeyFormat
	}

	s := base64.StdEncoding.EncodeToString(buf)

	return &PrivKey{
		Type: kt,
		S:    s,
		B:    []byte(s),
		K:    key,
	}, nil
}

// PubKey is a structure containing a public key in multiple formats
type PubKey struct {
	Type        KeyType     // Type of the the private key
	S           string      // String representation <type> <PEM key> <description>
	B           []byte      // Byte representation of string
	K           interface{} // Key interface{}
	Description string      // Optional description
}

// String converts a key to "<type> <key> <description>"
func (pk *PubKey) String() string {
	if pk == nil {
		return ""
	}

	s := fmt.Sprintf("%s %s %s", pk.Type.String(), pk.S, pk.Description)

	return strings.TrimSpace(s)
}

// Fingerprint return the fingerprint of the key
func (pk *PubKey) Fingerprint() string {
	binKey, _ := base64.StdEncoding.DecodeString(pk.S)
	b := sha256.Sum256(binKey)
	return hex.EncodeToString(b[:])
}

// MarshalJSON marshals a key into bytes
func (pk *PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pk.String())
}

// UnmarshalJSON unmarshals bytes into a key
func (pk *PubKey) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	pkCopy, err := PublicKeyFromString(s)
	if err != nil {
		return err
	}

	*pk = *pkCopy
	return err
}

// PublicKeyFromString creates a new public key based on the string data "<type> <key> <description>"
func PublicKeyFromString(data string) (*PubKey, error) {
	if !strings.Contains(data, " ") {
		return nil, errIncorrectKeyFormat
	}

	// <type> <data>
	parts := strings.SplitN(data, " ", 3)
	if len(parts) == 2 {
		parts = append(parts, "")
	}
	if len(parts) != 3 {
		return nil, errIncorrectKeyFormat
	}

	pk := &PubKey{}

	// Find the correct key type
	var err error
	pk.Type, err = FindKeyType(parts[0])
	if err != nil {
		return nil, errUnsupportedKeyType
	}

	// Set values
	pk.S = strings.TrimSpace(parts[1])
	pk.B = []byte(pk.S)
	pk.Description = parts[2]

	// Decode base64 before we parse to key
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(pk.B)))
	n, err := base64.StdEncoding.Decode(buf, pk.B)
	if err != nil {
		return nil, errIncorrectKeyFormat
	}

	// Parse key
	pk.K, err = pk.Type.ParsePublicKeyData(buf[:n])
	if err != nil {
		return nil, errIncorrectKeyFormat
	}

	return pk, nil
}

// PublicKeyFromInterface creates a new key based on an interface{} (like rsa.PublicKey)
func PublicKeyFromInterface(kt KeyType, key interface{}) (*PubKey, error) {
	buf, err := kt.ParsePublicKeyInterface(key)
	if err != nil {
		return nil, errIncorrectKeyFormat
	}

	s := base64.StdEncoding.EncodeToString(buf)

	return &PubKey{
		Type:        kt,
		S:           s,
		B:           []byte(s),
		K:           key,
		Description: "",
	}, nil
}
