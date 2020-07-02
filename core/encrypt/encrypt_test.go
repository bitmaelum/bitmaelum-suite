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
