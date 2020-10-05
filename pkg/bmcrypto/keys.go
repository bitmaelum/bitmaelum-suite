package bmcrypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

const (
	// KeyTypeRSA defines RSA keys
	KeyTypeRSA string = "rsa"
	// KeyTypeECDSA defines ECDSA keys
	KeyTypeECDSA string = "ecdsa"
	// KeyTypeED25519 defines ED25519 keys
	KeyTypeED25519 string = "ed25519"
)

// PrivKey is a structure containing a private key in multiple formats
type PrivKey struct {
	Type string      // Type of the private key
	S    string      // String representation <type> <PEM key>
	B    []byte      // Byte representation of string
	K    interface{} // Key interface{}
}

// PubKey is a structure containing a public key in multiple formats
type PubKey struct {
	Type        string      // Type of the the private key
	S           string      // String representation <type> <PEM key> <description>
	B           []byte      // Byte representation of string
	K           interface{} // Key interface{}
	Description string      // Optional description
}

// MarshalJSON marshals a key into bytes
func (pk *PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pk.String())
}

// String converts a key to "<type> <key> <description>"
func (pk *PubKey) String() string {
	return strings.TrimSpace(pk.Type + " " + pk.S + " " + pk.Description)
}

// UnmarshalJSON unmarshals bytes into a key
func (pk *PubKey) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if strings.HasPrefix(s, "-----") {
		keydata := strings.Split(s, "\n")
		s = "rsa " + strings.Join(keydata[1:len(keydata)-2], "")
	}

	pkCopy, err := NewPubKey(s)
	if err != nil {
		return err
	}

	// This seems wrong, but copy() doesn't work?
	pk.Type = pkCopy.Type
	pk.B = pkCopy.B
	pk.Description = pkCopy.Description
	pk.K = pkCopy.K
	pk.S = pkCopy.S

	return err
}

// MarshalJSON marshals a key into bytes
func (pk *PrivKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pk.String())
}

// Strings returns the key in a textual representation
func (pk *PrivKey) String() string {
	return pk.Type + " " + pk.S
}

// UnmarshalJSON unmarshals bytes into a key
func (pk *PrivKey) UnmarshalJSON(b []byte) error {
	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if strings.HasPrefix(s, "-----") {
		keydata := strings.Split(s, "\n")
		s = "rsa " + strings.Join(keydata[1:len(keydata)-2], "")
	}

	pkCopy, err := NewPrivKey(s)
	if err != nil {
		return err
	}

	// This seems wrong, but copy() doesn't work?
	pk.Type = pkCopy.Type
	pk.B = pkCopy.B
	pk.K = pkCopy.K
	pk.S = pkCopy.S

	return err
}

// CanEncrypt returns true if the key(type) is able to be used for encryption/decryption
func (pk *PubKey) CanEncrypt() bool {
	return pk.Type == KeyTypeRSA
}

// CanEncrypt returns true if the key(type) is able to be used for encryption/decryption
func (pk *PrivKey) CanEncrypt() bool {
	return pk.Type == KeyTypeRSA
}

// CanKeyExchange returns true if the key(type) is able to be used for key exchange
func (pk *PrivKey) CanKeyExchange() bool {
	if pk.Type == KeyTypeED25519 || pk.Type == KeyTypeECDSA {
		return true
	}
	return false
}

// CanKeyExchange returns true if the key(type) is able to be used for key exchange
func (pk *PubKey) CanKeyExchange() bool {
	if pk.Type == KeyTypeED25519 || pk.Type == KeyTypeECDSA {
		return true
	}
	return false
}

// NewPubKey creates a new public key based on the string data "<type> <key> <description>"
func NewPubKey(data string) (*PubKey, error) {
	pk := &PubKey{}

	if !strings.Contains(data, " ") {
		return nil, errors.New("incorrect key format")
	}

	// <type> <data> <description>
	parts := strings.SplitN(data, " ", 3)

	// Check and set type
	switch strings.ToLower(parts[0]) {
	case KeyTypeRSA:
		pk.Type = KeyTypeRSA
	case KeyTypeECDSA:
		pk.Type = KeyTypeECDSA
	case KeyTypeED25519:
		pk.Type = KeyTypeED25519
	default:
		return nil, errors.New("incorrect key type")
	}

	// Set values
	pk.S = strings.TrimSpace(parts[1])
	pk.B = []byte(pk.S)
	if len(parts) == 3 {
		pk.Description = parts[2]
	}

	// Decode base64 before we parse to key
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(pk.B)))
	n, err := base64.StdEncoding.Decode(buf, pk.B)
	if err != nil {
		return nil, errors.New("incorrect key data")
	}

	// Decode (base64-decoded) key
	pk.K, err = x509.ParsePKIXPublicKey(buf[:n])
	if err != nil {
		return nil, errors.New("incorrect key data")
	}

	return pk, nil
}

// NewPrivKey creates a new private key based on the string data "<type> <key>"
func NewPrivKey(data string) (*PrivKey, error) {
	pk := &PrivKey{}

	if !strings.Contains(data, " ") {
		return nil, errors.New("incorrect key format")
	}

	// <type> <data> <description>
	parts := strings.SplitN(data, " ", 2)

	// Check and set type
	switch strings.ToLower(parts[0]) {
	case KeyTypeRSA:
		pk.Type = KeyTypeRSA
	case KeyTypeECDSA:
		pk.Type = KeyTypeECDSA
	case KeyTypeED25519:
		pk.Type = KeyTypeED25519
	default:
		return nil, errors.New("incorrect key type")
	}

	// Set values
	pk.S = strings.TrimSpace(parts[1])
	pk.B = []byte(pk.S)

	// Decode base64 before we parse to key
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(pk.B)))
	n, err := base64.StdEncoding.Decode(buf, pk.B)
	if err != nil {
		return nil, errors.New("incorrect key data")
	}

	// Parse key, see which parses actually works..
	if pk.Type == KeyTypeECDSA {
		pk.K, err = x509.ParseECPrivateKey(buf[:n])
	} else {
		pk.K, err = x509.ParsePKCS1PrivateKey(buf[:n])
	}
	if err != nil {
		pk.K, err = x509.ParsePKCS8PrivateKey(buf[:n])
	}

	if err != nil {
		return nil, errors.New("incorrect key data")
	}

	return pk, nil
}

// NewPrivKeyFromInterface creates a new key based on an interface{} (like rsa.PrivateKey)
func NewPrivKeyFromInterface(key interface{}) (*PrivKey, error) {
	var t string
	switch key.(type) {
	case *rsa.PrivateKey:
		t = KeyTypeRSA
	case *ecdsa.PrivateKey:
		t = KeyTypeECDSA
	case ed25519.PrivateKey:
		t = KeyTypeED25519
	default:
		return nil, errors.New("unknown key type")
	}

	pk := &PrivKey{
		Type: t,
		K:    key,
	}

	buf, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, errors.New("incorrect key")
	}
	pk.S = base64.StdEncoding.EncodeToString(buf)
	pk.B = []byte(pk.S)

	return pk, nil
}

// NewPubKeyFromInterface creates a new key based on an interface{} (like rsa.PublicKey)
func NewPubKeyFromInterface(key interface{}) (*PubKey, error) {
	var t string
	switch key.(type) {
	case *rsa.PublicKey:
		t = KeyTypeRSA
	case *ecdsa.PublicKey:
		t = KeyTypeECDSA
	case ed25519.PublicKey:
		t = KeyTypeED25519
	default:
		return nil, errors.New("unknown key type")
	}

	pk := &PubKey{
		Type: t,
		K:    key,
	}

	buf, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, errors.New("incorrect key")
	}
	pk.S = base64.StdEncoding.EncodeToString(buf)
	pk.B = []byte(pk.S)

	return pk, nil
}
