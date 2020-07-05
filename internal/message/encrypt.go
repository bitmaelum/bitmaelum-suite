package message

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

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
