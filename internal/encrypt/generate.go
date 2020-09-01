package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"io"
)

var curveFunc = elliptic.P384
var rsaBits int = 2048

// Allows for easy mocking
var randReader io.Reader = rand.Reader

// GenerateKeyPair generates a private/public keypair based on the given type
func GenerateKeyPair(kt string) (*bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
	switch kt {
	case bmcrypto.KeyTypeRSA:
		return generateKeyPairRSA()
	case bmcrypto.KeyTypeECDSA:
		return generateKeyPairECDSA()
	case bmcrypto.KeyTypeED25519:
		return generateKeyPairED25519()
	}

	return nil, nil, errors.New("incorrect key type specified")
}

func generateKeyPairRSA() (*bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
	privRSAKey, err := rsa.GenerateKey(randReader, rsaBits)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := bmcrypto.NewPrivKeyFromInterface(privRSAKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := bmcrypto.NewPubKeyFromInterface(privKey.K.(*rsa.PrivateKey).Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

func generateKeyPairECDSA() (*bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
	privECDSAKey, err := ecdsa.GenerateKey(curveFunc(), randReader)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := bmcrypto.NewPrivKeyFromInterface(privECDSAKey)
	if err != nil {
		return nil, nil, err
	}
	pubKey, err := bmcrypto.NewPubKeyFromInterface(privKey.K.(*ecdsa.PrivateKey).Public())
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

func generateKeyPairED25519() (*bmcrypto.PrivKey, *bmcrypto.PubKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(randReader)
	if err != nil {
		return nil, nil, err
	}

	priv, err := bmcrypto.NewPrivKeyFromInterface(privKey)
	if err != nil {
		return nil, nil, err
	}
	pub, err := bmcrypto.NewPubKeyFromInterface(pubKey)
	if err != nil {
		return nil, nil, err
	}
	return priv, pub, nil
}
