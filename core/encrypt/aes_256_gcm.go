package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"github.com/bitmaelum/bitmaelum-server/internal/message"
	"io"
)

// JSONEncrypt encrypts a structure that is marshalled to JSON
func JSONEncrypt(key []byte, data interface{}) ([]byte, error) {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return MessageEncrypt(key, plaintext)
}

// JSONDecrypt decrypts data back from a encrypted marshalled JSON structure
func JSONDecrypt(key []byte, ciphertext []byte, v interface{}) error {
	plaintext, err := MessageDecrypt(key, ciphertext)
	if err != nil {
		return err
	}

	return json.Unmarshal(plaintext, &v)
}

// MessageEncrypt encrypts a binary message
func MessageEncrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := nonceGenerator(aead.NonceSize())
	if err != nil {
		return nil, err
	}

	return append(nonce, aead.Seal(nil, nonce, plaintext, nil)...), nil
}

// MessageDecrypt decrypts a binary message
func MessageDecrypt(key []byte, message []byte) ([]byte, error) {
	// Key should be 32byte (256bit)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()
	if len(message) < nonceSize {
		return nil, err
	}

	nonce, ciphertext := message[:nonceSize], message[nonceSize:]
	return aead.Open(nil, nonce, ciphertext, nil)
}

// CatalogEncrypt encrypts a catalog with a random key.
func CatalogEncrypt(catalog message.Catalog) ([]byte, []byte, error) {
	catalogKey, err := keyGenerator()
	if err != nil {
		return nil, nil, err
	}

	ciphertext, err := JSONEncrypt(catalogKey, catalog)
	if err != nil {
		return nil, nil, err
	}

	return catalogKey, ciphertext, nil
}

// CatalogDecrypt decrypts a catalog with the given key
func CatalogDecrypt(key, data []byte) (*message.Catalog, error) {
	// data, err := encode.Decode(data)
	// if err != nil {
	// 	return nil, err
	// }

	catalog := &message.Catalog{}
	err := JSONDecrypt(key, data, &catalog)
	if err != nil {
		return nil, err
	}

	return catalog, nil
}

// Generator to generate nonces for AEAD. Used so we can easily mock it in tests
var nonceGenerator = func(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, nonce)

	return nonce, err
}

// Generator to generate keys for catalog encryption. Used so we can easily mock it in tests
var keyGenerator = func() ([]byte, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)

	return key, err
}
