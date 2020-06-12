package encrypt

import (
    "crypto"
    "crypto/ecdsa"
    "crypto/ed25519"
    "crypto/rand"
    "crypto/rsa"
    "errors"
)

// Encrypt a message with the given key
func Encrypt(key interface{}, message []byte) ([]byte, error) {
    switch key.(type) {
	case *rsa.PrivateKey:
	    return encryptRsa(key.(*rsa.PrivateKey), message)
	case *ecdsa.PrivateKey:
		return encryptEcdsa(key.(*ecdsa.PrivateKey), message)
	case ed25519.PrivateKey:
		return encryptEd25519(key.(ed25519.PrivateKey), message)
	}

	return nil, errors.New("Unknown key type for signing")
}

func encryptEd25519(key ed25519.PrivateKey, message []byte) ([]byte, error) {
    return ed25519.Sign(key, message), nil
}

func encryptEcdsa(key *ecdsa.PrivateKey, message []byte) ([]byte, error) {
    return key.Sign(rand.Reader, message, nil)
}

func encryptRsa(key *rsa.PrivateKey, message []byte) ([]byte, error) {
    return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, message)
}
