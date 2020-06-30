package encrypt

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// PEMToPrivKey Converts a PEM & PKCS8 encoded private key
func PEMToPrivKey(pemData []byte) (interface{}, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("PEM decoding failed")
	}

	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

// PEMToPubKey Converts a PEM & PKCS8 encoded public key
func PEMToPubKey(pemData []byte) (interface{}, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("PEM decoding failed")
	}

	return x509.ParsePKIXPublicKey(block.Bytes)
}

// PrivKeyToPEM Convert a private key into PKCS8/PEM format
func PrivKeyToPEM(key interface{}) (string, error) {
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = pem.Encode(&b, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	return b.String(), err
}

// PubKeyToPEM Convert a public key into PKCS8/PEM format
func PubKeyToPEM(key interface{}) (string, error) {
	pubBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = pem.Encode(&b, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes})
	return b.String(), err
}
