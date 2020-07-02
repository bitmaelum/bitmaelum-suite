package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	pubPEM, privPEM, err := GenerateKeyPair(KeyTypeRSA)
	assert.Nil(t, err)
	pubKey, _ := PEMToPubKey([]byte(pubPEM))
	privKey, _ := PEMToPrivKey([]byte(privPEM))
	assert.IsType(t, (*rsa.PrivateKey)(nil), privKey)
	assert.IsType(t, (*rsa.PublicKey)(nil), pubKey)

	pubPEM, privPEM, err = GenerateKeyPair(KeyTypeECDSA)
	assert.Nil(t, err)
	pubKey, _ = PEMToPubKey([]byte(pubPEM))
	privKey, _ = PEMToPrivKey([]byte(privPEM))
	assert.IsType(t, (*ecdsa.PrivateKey)(nil), privKey)
	assert.IsType(t, (*ecdsa.PublicKey)(nil), pubKey)

	pubPEM, privPEM, err = GenerateKeyPair(KeyTypeED25519)
	assert.Nil(t, err)
	pubKey, _ = PEMToPubKey([]byte(pubPEM))
	privKey, _ = PEMToPrivKey([]byte(privPEM))
	assert.IsType(t, (ed25519.PrivateKey)(nil), privKey)
	assert.IsType(t, (ed25519.PublicKey)(nil), pubKey)

	pubPEM, privPEM, err = GenerateKeyPair(25)
	assert.EqualError(t, err, "incorrect key type specified")
}
