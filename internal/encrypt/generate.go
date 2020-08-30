package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"errors"
)

var curveFunc = elliptic.P384
var rsaBits int = 2048

// GenerateKeyPair generates a private/public keypair based on the given type
func GenerateKeyPair(kt string) (*PrivKey, *PubKey, error) {
	switch kt {
	case KeyTypeRSA:
		return generateKeyPairRSA()
	case KeyTypeECDSA:
		return generateKeyPairECDSA()
	case KeyTypeED25519:
		return generateKeyPairED25519()
	}

	return nil, nil, errors.New("incorrect key type specified")
}

func generateKeyPairRSA() (*PrivKey, *PubKey, error) {
	privRSAKey, err := rsa.GenerateKey(randReader, rsaBits)
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

func generateKeyPairECDSA() (*PrivKey, *PubKey, error) {
	privECDSAKey, err := ecdsa.GenerateKey(curveFunc(), randReader)
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

func generateKeyPairED25519() (*PrivKey, *PubKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(randReader)
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
