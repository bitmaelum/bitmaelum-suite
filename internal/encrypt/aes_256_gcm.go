package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
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
func CatalogEncrypt(catalog interface{}) ([]byte, []byte, error) {
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
func CatalogDecrypt(key, data []byte, catalog interface{}) error {
	err := JSONDecrypt(key, data, catalog)
	if err != nil {
		return err
	}

	return nil
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

// GenerateIvAndKey generate a random IV and key
// @TODO: This is not a good spot. We should store it in the encrypt page, but this gives us a import cycle
func GenerateIvAndKey() ([]byte, []byte, error) {
	iv := make([]byte, 16)
	n, err := rand.Read(iv)
	if n != 16 || err != nil {
		return nil, nil, err
	}

	key := make([]byte, 32)
	n, err = rand.Read(key)
	if n != 32 || err != nil {
		return nil, nil, err
	}

	return iv, key, nil
}

// GetAesEncryptorReader returns a reader that automatically encrypts reader blocks through CFB stream
// @TODO: This is not a good spot. We should store it in the encrypt page, but this gives us a import cycle
func GetAesEncryptorReader(iv []byte, key []byte, r io.Reader) (io.Reader, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, err
}

// GetAesDecryptorReader returns a reader that automatically decrypts reader blocks through CFB stream
func GetAesDecryptorReader(iv []byte, key []byte, r io.Reader) (io.Reader, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, err
}
