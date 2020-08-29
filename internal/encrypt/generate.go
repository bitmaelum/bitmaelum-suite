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

var curveFunc = elliptic.P384
var rsaBits int = 2048

// GenerateKeyPair generates a private/public keypair based on the given type
func GenerateKeyPair(kt string) (*PrivKey, *PubKey, error) {
	if kt == KeyTypeRSA {

		privRSAKey, err := rsa.GenerateKey(rand.Reader, rsaBits)
		if err != nil {
			return nil, nil, err
		}

		privKey, err := NewPrivKeyFromInterface(privRSAKey)
		if err != nil {
			return nil, nil, err
		}
		pubKey, err := NewPubKeyFromInterface(privKey.K.(*rsa.PrivateKey).Public())
		if err != nil {
			return nil, nil, err
		}

		return privKey, pubKey, nil
	}

	if kt == KeyTypeECDSA {
		var err error
		privECDSAKey, err := ecdsa.GenerateKey(curveFunc(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}

		privKey, err := NewPrivKeyFromInterface(privECDSAKey)
		if err != nil {
			return nil, nil, err
		}
		pubKey, err := NewPubKeyFromInterface(privKey.K.(*ecdsa.PrivateKey).Public())
		if err != nil {
			return nil, nil, err
		}

		return privKey, pubKey, nil
	}

	if kt == KeyTypeED25519 {
		var err error
		pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, err
		}

		priv, err := NewPrivKeyFromInterface(privKey)
		if err != nil {
			return nil, nil, err
		}
		pub, err := NewPubKeyFromInterface(pubKey)
		if err != nil {
			return nil, nil, err
		}
		return priv, pub, nil
	}

	return nil, nil, errors.New("incorrect key type specified")
}
