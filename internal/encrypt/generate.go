package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

type keyType int

const (
	// KeyTypeRSA RSA key
	KeyTypeRSA = iota
	// KeyTypeECDSA Elliptic curve key
	KeyTypeECDSA
	// KeyTypeED25519 ED25519 key
	KeyTypeED25519
)

var curveFunc = elliptic.P384
var rsaBits int = 2048

// GenerateKeyPair generates a public/private keypair based on the given type
func GenerateKeyPair(kt keyType) (string, string, error) {
	var privKey, pubKey interface{}

	switch kt {
	case KeyTypeRSA:
		var err error
		privKey, err = rsa.GenerateKey(rand.Reader, rsaBits)
		if err != nil {
			return "", "", err
		}
		pubKey = privKey.(*rsa.PrivateKey).Public()
	case KeyTypeECDSA:
		var err error
		privKey, err = ecdsa.GenerateKey(curveFunc(), rand.Reader)
		if err != nil {
			return "", "", err
		}
		pubKey = privKey.(*ecdsa.PrivateKey).Public()
	case KeyTypeED25519:
		var err error
		pubKey, privKey, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", errors.New("incorrect key type specified")
	}

	privKeyPem, err := PrivKeyToPEM(privKey)
	if err != nil {
		return "", "", err
	}
	pubKeyPem, err := PubKeyToPEM(pubKey)
	if err != nil {
		return "", "", err
	}

	return pubKeyPem, privKeyPem, nil
}