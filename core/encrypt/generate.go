package encrypt

import (
    "crypto/ecdsa"
    "crypto/ed25519"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/rsa"
    "errors"
)

type KeyType int

const (
	KeyTypeRSA = iota
	KeyTypeECDSA
	KeyTypeED25519
)

var curveFunc = elliptic.P384
var rsa_bits int = 2048

type JustAKey = interface{}

// Generates a public/private keypair based on the given type
func GenerateKeyPair(kt KeyType) (string, string, error) {
    var privKey, pubKey JustAKey

	switch kt {
    case KeyTypeRSA:
        var err error
        privKey, err = rsa.GenerateKey(rand.Reader, rsa_bits)
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
        privKey, pubKey, err = ed25519.GenerateKey(rand.Reader)
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
