package bmcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionId(t *testing.T) {
	txID := &TransactionID{
		P: []byte{120, 216, 242, 172, 92, 92, 168, 227, 245, 83, 247, 191, 240, 109, 136, 59, 230, 226, 60, 74, 192, 193, 188, 164, 195, 112, 46, 42, 177, 238, 134, 210},
		R: []byte{4, 173, 159, 128, 130, 144, 107, 116, 74, 189, 217, 50, 76, 127, 250, 119, 30, 51, 208, 135, 247, 127, 92, 44, 255, 60, 131, 211, 92, 187, 57, 65},
	}
	assert.Equal(t, "78d8f2ac5c5ca8e3f553f7bff06d883be6e23c4ac0c1bca4c3702e2ab1ee86d204ad9f8082906b744abdd9324c7ffa771e33d087f77f5c2cff3c83d35cbb3941", txID.ToHex())

	txID, err := TxIDFromString("78d8f2ac5c5ca8e3f553f7bff06d883be6e23c4ac0c1bca4c3702e2ab1ee86d204ad9f8082906b744abdd9324c7ffa771e33d087f77f5c2cff3c83d35cbb3941")
	assert.NoError(t, err)
	assert.Equal(t, []byte{120, 216, 242, 172, 92, 92, 168, 227, 245, 83, 247, 191, 240, 109, 136, 59, 230, 226, 60, 74, 192, 193, 188, 164, 195, 112, 46, 42, 177, 238, 134, 210}, txID.P)
	assert.Equal(t, []byte{4, 173, 159, 128, 130, 144, 107, 116, 74, 189, 217, 50, 76, 127, 250, 119, 30, 51, 208, 135, 247, 127, 92, 44, 255, 60, 131, 211, 92, 187, 57, 65}, txID.R)

	txID, err = TxIDFromString("78d8f2ac5c5ca8e3f553f7bff06d883be6e23c4ac0c1bca4c3702e2ab1ee86d204ad9f8082906b744a")
	assert.Error(t, err)
	assert.Nil(t, txID)

	txID, err = TxIDFromString("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX7bff06d883be6e23c4ac0c1bca4c3702e2ab1ee86d204ad9f8082906b744a")
	assert.Error(t, err)
	assert.Nil(t, txID)
}

func TestWrongKeyType(t *testing.T) {
	privKey, pubKey, _ := generateKeyPairRSA()

	_, _, err := DualKeyExchange(*pubKey)
	assert.Error(t, err)

	_, _, err = DualKeyGetSecret(*privKey, TransactionID{})
	assert.Error(t, err)
}

func TestDualKeyExchange(t *testing.T) {

	privKey, _ := NewPrivKey("ed25519 MC4CAQAwBQYDK2VwBCIEIBJsN8lECIdeMHEOZhrdDNEZl5BuULetZsbbdsZBjZ8a")
	pubKey, _ := NewPubKey("ed25519 MCowBQYDK2VwAyEAblFzZuzz1vItSqdHbr/3DZMYvdoy17ALrjq3BM7kyKE=")

	D, txID, err := DualKeyExchange(*pubKey)
	assert.NoError(t, err)

	Dprime, ok, err := DualKeyGetSecret(*privKey, *txID)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, D, Dprime)

	privKey2, _ := NewPrivKey("ed25519 MC4CAQAwBQYDK2VwBCIEII6nA1nsVQu1Pid+CoH6yxw9Z2yOU9++S30awQvIW3m/")
	Dprime, ok, err = DualKeyGetSecret(*privKey2, *txID)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, Dprime)
}

func BenchmarkDualKeyExchange(b *testing.B) {
	privKey, _ := NewPrivKey("ed25519 MC4CAQAwBQYDK2VwBCIEIBJsN8lECIdeMHEOZhrdDNEZl5BuULetZsbbdsZBjZ8a")
	pubKey, _ := NewPubKey("ed25519 MCowBQYDK2VwAyEAblFzZuzz1vItSqdHbr/3DZMYvdoy17ALrjq3BM7kyKE=")

	for i := 0; i < b.N; i++ {
		D, txID, err := DualKeyExchange(*pubKey)
		assert.NoError(b, err)

		Dprime, ok, err := DualKeyGetSecret(*privKey, *txID)
		assert.NoError(b, err)
		assert.True(b, ok)
		assert.Equal(b, D, Dprime)
	}
}
