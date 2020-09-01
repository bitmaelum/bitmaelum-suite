package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	message = []byte("this is the message we need to sign")
)

func TestGenerate(t *testing.T) {
	privKey, pubKey, err := GenerateKeyPair(bmcrypto.KeyTypeRSA)
	assert.Nil(t, err)
	assert.IsType(t, (*rsa.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (*rsa.PublicKey)(nil), pubKey.K)

	// Check if we can verify with this key
	sig, err := bmcrypto.Sign(*privKey, message)
	assert.Nil(t, err)
	b, err := bmcrypto.Verify(*pubKey, message, sig)
	assert.Nil(t, err)
	assert.True(t, b)

	privKey, pubKey, err = GenerateKeyPair(bmcrypto.KeyTypeECDSA)
	assert.Nil(t, err)
	assert.IsType(t, (*ecdsa.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (*ecdsa.PublicKey)(nil), pubKey.K)

	// Check if we can verify with this key
	sig, err = bmcrypto.Sign(*privKey, message)
	assert.Nil(t, err)
	b, err = bmcrypto.Verify(*pubKey, message, sig)
	assert.Nil(t, err)
	assert.True(t, b)

	privKey, pubKey, err = GenerateKeyPair(bmcrypto.KeyTypeED25519)
	assert.Nil(t, err)
	assert.IsType(t, (ed25519.PrivateKey)(nil), privKey.K)
	assert.IsType(t, (ed25519.PublicKey)(nil), pubKey.K)

	// Check if we can verify with this key
	sig, err = bmcrypto.Sign(*privKey, message)
	assert.Nil(t, err)
	b, err = bmcrypto.Verify(*pubKey, message, sig)
	assert.Nil(t, err)
	assert.True(t, b)

	// Unknown key
	_, _, err = GenerateKeyPair("foobar")
	assert.EqualError(t, err, "incorrect key type specified")
}
