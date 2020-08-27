package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"strings"
)

const (
	KEYTYPE_RSA     string = "rsa"
	KEYTYPE_ECDSA   string = "ecdsa"
	KEYTYPE_ED25519 string = "ed25519"
)

type PrivKey struct {
	Type        string
	S           string
	B           []byte
	K           interface{}
}

type PubKey struct {
	Type string
	S    string
	B    []byte
	K    interface{}
	Description string
}

func NewPubKey(data string) (*PubKey, error) {
	pk := &PubKey{}

	if !strings.Contains(data, " ") {
		return nil, errors.New("incorrect key format")
	}

	// <type> <data> <description>
	parts := strings.SplitN(data, " ", 3)

	// Check and set type
	switch strings.ToLower(parts[0]) {
	case KEYTYPE_RSA:
		pk.Type = KEYTYPE_RSA
	case KEYTYPE_ECDSA:
		pk.Type = KEYTYPE_ECDSA
	case KEYTYPE_ED25519:
		pk.Type = KEYTYPE_ED25519
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

func NewPrivKey(data string) (*PrivKey, error) {
	pk := &PrivKey{}

	if !strings.Contains(data, " ") {
		return nil, errors.New("incorrect key format")
	}

	// <type> <data> <description>
	parts := strings.SplitN(data, " ", 2)

	// Check and set type
	switch strings.ToLower(parts[0]) {
	case KEYTYPE_RSA:
		pk.Type = KEYTYPE_RSA
	case KEYTYPE_ECDSA:
		pk.Type = KEYTYPE_ECDSA
	case KEYTYPE_ED25519:
		pk.Type = KEYTYPE_ED25519
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

	// Decode (base64-decoded) key
	if pk.Type == KEYTYPE_ECDSA {
		pk.K, err = x509.ParseECPrivateKey(buf[:n])
	} else {
		pk.K, err = x509.ParsePKCS8PrivateKey(buf[:n])
	}
	if err != nil {
		return nil, errors.New("incorrect key data")
	}

	return pk, nil
}


func NewPrivKeyFromInterface(key interface{}) (*PrivKey, error) {
	switch key.(type) {
	case *rsa.PrivateKey:
		return nil, nil
	case *ecdsa.PrivateKey:
		return nil, nil
	case ed25519.PrivateKey:
		return nil, nil
	}

	return nil, errors.New("incorrect key type")
}

func NewPubKeyFromInterface(key interface{}) (*PubKey, error) {
	switch key.(type) {
	case *rsa.PublicKey:
		return nil, nil
	case *ecdsa.PublicKey:
		return nil, nil
	case ed25519.PublicKey:
		return nil, nil
	}

	return nil, errors.New("incorrect key type")
}
