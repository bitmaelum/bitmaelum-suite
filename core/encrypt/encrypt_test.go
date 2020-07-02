package encrypt

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestEncrypt(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/mykey.pub")
	pubKey, _ := PEMToPubKey(data)

	data, _ = ioutil.ReadFile("./testdata/mykey.pem")
	privKey, _ := PEMToPrivKey(data)

	cipher, err := Encrypt(pubKey, []byte("foobar"))
	assert.Nil(t, err)
	assert.NotEqual(t, []byte("foobar"), cipher)

	plaintext, err := Decrypt(privKey, cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)
}

func TestEncryptWithIncorrectKeys(t *testing.T) {
	var err error
	
	_, err = Encrypt("incorrectkeyinterface", []byte("foobar"))
	assert.EqualError(t, err, "Unknown key type for signing")

	_, err = Decrypt("incorrectkeyinterface", []byte("foobar"))
	assert.EqualError(t, err, "Unknown key type for signing")

}


/*
// Encrypt a message with the given key
func Encrypt(key interface{}, message []byte) ([]byte, error) {
	switch key.(type) {
	case *rsa.PublicKey:
		return encryptRsa(key.(*rsa.PublicKey), message)
	case *ecdsa.PublicKey:
	case ed25519.PublicKey:
		return nil, errors.New("This key type is not usable for encryption")
	}

	return nil, errors.New("Unknown key type for signing")
}

// Decrypt a message with the given key
func Decrypt(key interface{}, message []byte) ([]byte, error) {
	switch key.(type) {
	case *rsa.PrivateKey:
		return decryptRsa(key.(*rsa.PrivateKey), message)
	case *ecdsa.PrivateKey:
	case ed25519.PrivateKey:
		return nil, errors.New("This key type is not usable for encryption")
	}

	return nil, errors.New("Unknown key type for signing")
}

func encryptRsa(key *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, key, message)
}

func decryptRsa(key *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, message)
}
*/
