package bmcrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
)

// Encrypt a message with the given key
func Encrypt(key PubKey, message []byte) ([]byte, error) {
	if !key.CanEncrypt() {
		return nil, errors.New("This key type is not usable for encryption")
	}

	switch key.Type {
	case KeyTypeRSA:
		return encryptRsa(key.K.(*rsa.PublicKey), message)
	}

	return nil, errors.New("encryption not implemented for" + key.Type)
}

// Decrypt a message with the given key
func Decrypt(key PrivKey, message []byte) ([]byte, error) {
	if !key.CanEncrypt() {
		return nil, errors.New("This key type is not usable for encryption")
	}

	switch key.Type {
	case KeyTypeRSA:
		return decryptRsa(key.K.(*rsa.PrivateKey), message)
	}

	return nil, errors.New("encryption not implemented for" + key.Type)
}

func encryptRsa(key *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, key, message)
}

func decryptRsa(key *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, message)
}
