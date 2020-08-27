package encrypt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var TestKeySet1 = []string{
	"rsa MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8pvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5aG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxXPOhC1QUcllt08yBB3wIDAQAB",
	"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8pvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5aG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxXPOhC1QUcllt08yBB3wIDAQAB",
	"fooobar",
	"rsa foo",
	"foo MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8pvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5aG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxXPOhC1QUcllt08yBB3wIDAQAB",
}
var TestKeySet2 = []string{
	"ed25519 MCowBQYDK2VwAyEAOT5K9VqRFjJ2Q0RWWqT8+OndKQrmzXCuMTSnhfvVpws= my foo bar",
}
var TestKeySet3 = []string{
	"ecdsa MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE58d01mIg3iWGqBHY7N6ch4L4LWya8Es2luWC08Wjn994nLgIOUp+cdMUfDYBe/x1aPE/yghGe3rrF4jW8uxVWy40BZK4NIu5yjMgSw0WBGTxOmZsVaA/xaOzvZSMTXxM",
}
var TestKeySet4 = []string{
	"ecdsa MIGkAgEBBDBLD4tDPxb/Xw2SzOsDEwl42LinqQmlWmcusiQJSnHn2VJsHzTuBoj7zE0dGhBS/ESgBwYFK4EEACKhZANiAATnx3TWYiDeJYaoEdjs3pyHgvgtbJrwSzaW5YLTxaOf33icuAg5Sn5x0xR8NgF7/HVo8T/KCEZ7eusXiNby7FVbLjQFkrg0i7nKMyBLDRYEZPE6ZmxVoD/Fo7O9lIxNfEw=",
}
var TestKeySet5 = []string{
	"ed25519 MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7",
}




func TestRSAPubKey(t *testing.T) {
	// Correct
	pk, err := NewPubKey(TestKeySet1[0])
	assert.NoError(t, err)
	assert.Equal(t, "rsa", pk.Type)
	assert.Equal(t, "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8pvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5aG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxXPOhC1QUcllt08yBB3wIDAQAB", pk.S)
	assert.Equal(t, 65537, pk.K.(*rsa.PublicKey).E)
	assert.Equal(t, "", pk.Description)

	// Without type
	pk, err = NewPubKey(TestKeySet1[1])
	assert.EqualError(t, err, "incorrect key format")

	// Incorrect too
	pk, err = NewPubKey(TestKeySet1[2])
	assert.EqualError(t, err, "incorrect key format")

	// right type, wrong data
	pk, err = NewPubKey(TestKeySet1[3])
	assert.EqualError(t, err, "incorrect key data")

	// wrong type, right data
	pk, err = NewPubKey(TestKeySet1[4])
	assert.EqualError(t, err, "incorrect key type")
}

func TestED25519PubKey(t *testing.T) {
	key := ed25519.PublicKey{0x39, 0x3e, 0x4a, 0xf5, 0x5a, 0x91, 0x16, 0x32, 0x76, 0x43, 0x44, 0x56, 0x5a, 0xa4, 0xfc, 0xf8, 0xe9, 0xdd, 0x29, 0xa, 0xe6, 0xcd, 0x70, 0xae, 0x31, 0x34, 0xa7, 0x85, 0xfb, 0xd5, 0xa7, 0xb}
	// Correct ED25519
	pk, err := NewPubKey(TestKeySet2[0])
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", pk.Type)
	assert.Equal(t, "MCowBQYDK2VwAyEAOT5K9VqRFjJ2Q0RWWqT8+OndKQrmzXCuMTSnhfvVpws=", pk.S)
	assert.Equal(t, key, pk.K.(ed25519.PublicKey))
	assert.Equal(t, "my foo bar", pk.Description)
}

func TestECDSAPub(t *testing.T) {
	x := new(big.Int)
	x.SetString("35674072579805598781611597844361315596081721500728226840731099278229644034819246176780429354491654135480285061184629", 10)
	y := new(big.Int)
	y.SetString("16152110512098221797044896547068689859301830683829073190715770376238845227548146305554820081396880025274135358569548", 10)

	pk, err := NewPubKey(TestKeySet3[0])
	assert.NoError(t, err)
	assert.Equal(t, "ecdsa", pk.Type)
	assert.Equal(t, "MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE58d01mIg3iWGqBHY7N6ch4L4LWya8Es2luWC08Wjn994nLgIOUp+cdMUfDYBe/x1aPE/yghGe3rrF4jW8uxVWy40BZK4NIu5yjMgSw0WBGTxOmZsVaA/xaOzvZSMTXxM", pk.S)
	assert.Equal(t, x, pk.K.(*ecdsa.PublicKey).X)
	assert.Equal(t, y, pk.K.(*ecdsa.PublicKey).Y)
	assert.Equal(t, "P-384", pk.K.(*ecdsa.PublicKey).Curve.Params().Name)
	assert.Equal(t, "", pk.Description)

	// Check private key data
	pk, err = NewPubKey(TestKeySet4[0])
	assert.EqualError(t, err, "incorrect key data")
}

func TestECDSAPriv(t *testing.T) {
	d := new(big.Int)
	d.SetString("11552901970705313238876759535655836311969175439875617508331015348976563282371780402868289145078300182831985971100740", 10)

	pk, err := NewPrivKey(TestKeySet4[0])
	assert.NoError(t, err)
	assert.Equal(t, "ecdsa", pk.Type)
	assert.Equal(t, "MIGkAgEBBDBLD4tDPxb/Xw2SzOsDEwl42LinqQmlWmcusiQJSnHn2VJsHzTuBoj7zE0dGhBS/ESgBwYFK4EEACKhZANiAATnx3TWYiDeJYaoEdjs3pyHgvgtbJrwSzaW5YLTxaOf33icuAg5Sn5x0xR8NgF7/HVo8T/KCEZ7eusXiNby7FVbLjQFkrg0i7nKMyBLDRYEZPE6ZmxVoD/Fo7O9lIxNfEw=", pk.S)
	assert.Equal(t, d, pk.K.(*ecdsa.PrivateKey).D)
	assert.Equal(t, "P-384", pk.K.(*ecdsa.PrivateKey).Curve.Params().Name)
}

func TestED25519Priv(t *testing.T) {
	b := []byte{0xc7, 0x54, 0xb9, 0x99, 0x3c, 0x19, 0x73, 0xc, 0x46, 0x7a, 0xdb, 0x67, 0xdd, 0x9, 0xc, 0xd2, 0x2b, 0xcf, 0x98, 0x7b, 0x69, 0xcd, 0xb0, 0x85, 0x98, 0xaa, 0xa6, 0xf7, 0x83, 0x4a, 0x6, 0x7b}

	pk, err := NewPrivKey(TestKeySet5[0])
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", pk.Type)
	assert.Equal(t, "MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7", pk.S)
	assert.Equal(t, b, pk.K.(ed25519.PrivateKey).Seed())
}
