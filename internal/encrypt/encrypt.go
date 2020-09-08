package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// Encrypt a message with the given key
func Encrypt(key bmcrypto.PubKey, message []byte) ([]byte, error) {
	if !key.CanEncrypt() {
		return nil, errors.New("This key type is not usable for encryption")
	}

	switch key.Type {
	case bmcrypto.KeyTypeRSA:
		return encryptRsa(key.K.(*rsa.PublicKey), message)
	}

	return nil, errors.New("encryption not implemented for" + key.Type)
}

// Decrypt a message with the given key
func Decrypt(key bmcrypto.PrivKey, message []byte) ([]byte, error) {
	if !key.CanEncrypt() {
		return nil, errors.New("This key type is not usable for encryption")
	}

	switch key.Type {
	case bmcrypto.KeyTypeRSA:
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
