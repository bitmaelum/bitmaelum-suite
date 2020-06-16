package encrypt

import (
    "crypto/x509"
    "encoding/pem"
    "errors"
)

// Converts a PEM & PKCS8 encoded private key
func PEMToPrivKey(pemData []byte) (interface{}, error) {
    block, _ := pem.Decode(pemData)
    if block == nil {
        return nil, errors.New("PEM decoding failed")
    }

    return x509.ParsePKCS8PrivateKey(block.Bytes)
}

// Converts a PEM & PKCS8 encoded public key
func PEMToPubKey(pemData []byte) (interface{}, error) {
    block, _ := pem.Decode(pemData)
    if block == nil {
        return nil, errors.New("PEM decoding failed")
    }

    return x509.ParsePKIXPublicKey(block.Bytes)
}

// Convert a private key into PKCS8/PEM format
func PrivKeyToPEM(key interface{}) ([]byte, error) {
    return x509.MarshalPKCS8PrivateKey(key)
}

// Convert a public key into PKCS8/PEM format
func PubKeyToPEM(key interface{}) ([]byte, error) {
    return x509.MarshalPKIXPublicKey(key)
}



