package internal

import (
	"crypto/ed25519"
	"testing"

	testing2 "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestSigningMethodEdDSAAlg(t *testing.T) {
	m := &SigningMethodEdDSA{}
	assert.Equal(t, "EdDSA", m.Alg())
}

func TestSigningMethodEdDSASign(t *testing.T) {
	m := &SigningMethodEdDSA{}

	privKey, pubKey, _ := testing2.ReadTestKey("../testdata/key-ed25519-1.json")

	s, err := m.Sign("foobar", privKey.K.(ed25519.PrivateKey))
	assert.NoError(t, err)

	err = m.Verify("foobar", s, pubKey.K.(ed25519.PublicKey))
	assert.NoError(t, err)

	err = m.Verify("foobarfoofofoo", s, pubKey.K.(ed25519.PublicKey))
	assert.Error(t, err)
}
