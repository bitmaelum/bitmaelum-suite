package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	privKey, pubKey, err := GenerateKeyPair(KeyTypeRSA)
	assert.Nil(t, err)
	assert.IsType(t, (*rsa.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (*rsa.PublicKey)(nil), pubKey.K)

	privKey, pubKey, err = GenerateKeyPair(KeyTypeECDSA)
	assert.Nil(t, err)
	assert.IsType(t, (*ecdsa.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (*ecdsa.PublicKey)(nil), pubKey.K)

	privKey, pubKey, err = GenerateKeyPair(KeyTypeED25519)
	assert.Nil(t, err)
	assert.IsType(t, (ed25519.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (ed25519.PublicKey)(nil), pubKey.K)

	_, _, err = GenerateKeyPair(25)
	assert.EqualError(t, err, "incorrect key type specified")
}
