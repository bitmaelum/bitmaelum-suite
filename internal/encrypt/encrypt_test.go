package encrypt

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestEncrypt(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/mykey.pub")
	pubKey, _ := NewPubKey(string(data))

	data, _ = ioutil.ReadFile("./testdata/mykey.pem")
	privKey, _ := NewPrivKey(string(data))

	cipher, err := Encrypt(*pubKey, []byte("foobar"))
	assert.Nil(t, err)
	assert.NotEqual(t, []byte("foobar"), cipher)

	plaintext, err := Decrypt(*privKey, cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)
}
