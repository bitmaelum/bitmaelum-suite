package bmcrypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"encoding/asn1"
	"errors"
	"io"
	"math/big"
)

// Allows for easy mocking
var randReader io.Reader = rand.Reader

type ecdsaSignature struct {
	R, S *big.Int
}

// Sign a message based on the given key.
func Sign(key PrivKey, message []byte) ([]byte, error) {
	switch key.Type {
	case KeyTypeRSA:
		h := crypto.SHA256.New()
		h.Write(message)
		hash := h.Sum(nil)

		return rsa.SignPKCS1v15(randReader, key.K.(*rsa.PrivateKey), crypto.SHA256, hash[:])
	case KeyTypeECDSA:
		r, s, err := ecdsa.Sign(randReader, key.K.(*ecdsa.PrivateKey), message)
		if err != nil {
			return nil, err
		}

		sig := ecdsaSignature{
			R: r,
			S: s,
		}

		return asn1.Marshal(sig)
	case KeyTypeED25519:
		return ed25519.Sign(key.K.(ed25519.PrivateKey), message), nil
	}

	return nil, errors.New("Unknown key type for signing")
}

// Verify if hash compares against the signature of the message
func Verify(key PubKey, message []byte, sig []byte) (bool, error) {
	switch key.Type {
	case KeyTypeRSA:
		h := crypto.SHA256.New()
		h.Write(message)
		hash := h.Sum(nil)

		err := rsa.VerifyPKCS1v15(key.K.(*rsa.PublicKey), crypto.SHA256, hash[:], sig)
		return err == nil, err
	case KeyTypeECDSA:
		ecdsaSig := ecdsaSignature{}
		_, err := asn1.Unmarshal(sig, &ecdsaSig)
		if err != nil {
			return false, err
		}

		return ecdsa.Verify(key.K.(*ecdsa.PublicKey), message, ecdsaSig.R, ecdsaSig.S), nil
	case KeyTypeED25519:
		return ed25519.Verify(key.K.(ed25519.PublicKey), message, sig), nil
	}

	return false, errors.New("Unknown key type for signing")
}
