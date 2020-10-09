package bmcrypto

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	//RSA Encryption
	data, _ := ioutil.ReadFile("../../testdata/pubkey.rsa")
	pubKey, _ := NewPubKey(string(data))

	data, _ = ioutil.ReadFile("../../testdata/privkey.rsa")
	privKey, _ := NewPrivKey(string(data))

	cipher, _, _, err := Encrypt(*pubKey, []byte("foobar"))
	assert.Nil(t, err)
	assert.NotEqual(t, []byte("foobar"), cipher)

	plaintext, err := Decrypt(*privKey, "", cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)

	//ED25519 Dual Key-Exchange + Encryption
	priv25519Key, _ := NewPrivKey("ed25519 MC4CAQAwBQYDK2VwBCIEIBJsN8lECIdeMHEOZhrdDNEZl5BuULetZsbbdsZBjZ8a")
	pub25519Key, _ := NewPubKey("ed25519 MCowBQYDK2VwAyEAblFzZuzz1vItSqdHbr/3DZMYvdoy17ALrjq3BM7kyKE=")
	cipher, txID, _, err := Encrypt(*pub25519Key, []byte("foobar"))
	assert.Nil(t, err)
	assert.NotEqual(t, []byte("foobar"), cipher)

	plaintext, err = Decrypt(*priv25519Key, txID, cipher)
	assert.Nil(t, err)
	assert.Equal(t, []byte("foobar"), plaintext)
}
