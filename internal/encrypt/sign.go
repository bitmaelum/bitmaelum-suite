package encrypt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"errors"
	"io"
)

// Allows for easy mocking
var randReader io.Reader = rand.Reader

// Sign a message based on the given key.
func Sign(key PrivKey, message []byte) ([]byte, error) {
	switch key.Type {
	case KEYTYPE_RSA:
		return rsa.SignPKCS1v15(randReader, key.K.(*rsa.PrivateKey), crypto.SHA256, message)
	case KEYTYPE_ECDSA:
		privkey := key.K.(*ecdsa.PrivateKey)
		return privkey.Sign(randReader, message, nil)
	case KEYTYPE_ED25519:
		return ed25519.Sign(key.K.(ed25519.PrivateKey), message), nil
	}

	return nil, errors.New("Unknown key type for signing")
}

// Verify if hash compares against the signature of the message
func Verify(key PrivKey, message []byte, hash []byte) (bool, error) {
	computedHash, err := Sign(key, message)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(computedHash, hash) == 1, nil
}
