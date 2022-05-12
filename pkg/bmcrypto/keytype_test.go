// Copyright (c) 2022 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package bmcrypto

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestKeySet1 = []string{
	"rsa MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAySB+eJGwb1Wu8KERTBLxKfMaZ7VRP92Q0LmJbqF5X2NrLyO+WkcYCcTqqoeW9unfoEiT4dS0I2YRLNcXztTaPs9xcoAWuXxPHwMKdjFKWQubttUoAAU08HpCkceO2y29NsrgOfKgUUc+D8FRXtyAe+GUg6wpWhxtT1BmjIgT4LIJKTKJ7cUdLxDQuKmT0uj0B+LQ4DTo4SvE7aV+Lb4wD2kwB0CHckHsbVC1jFPPcSLEtaSAjssaHe6v2c0oS2VdvuHq2pMa2lnuR/+wXeWn+zf0Rft0rcZNbIDpRt+DQH5g6VhA1Ww802mB/XYGaZcljzR5ArCDX8flvHP04ZcVWQIDAQAB",
	// "rsa MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8pvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5aG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxXPOhC1QUcllt08yBB3wIDAQAB",
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
	"ecdsa MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDDt8LNhN1AHVrBuBqKrryGgUS0EdRdjRuYGjONDYI96Vy8IcGsz4HQrX0biOzfE5iuhZANiAAT67cjQyt3qJUktq0dJy/KZ/15NhPpqBlG7NCwVyeyqcU2IlpE0bM+58BOkBpHCYq7zxEfXurDYIuCMKNKExJXXUeCczPBNHg9pWVVCPTHkypb69VURfYgQSo/58vSXCjU=",
	"ecdsa MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE+u3I0Mrd6iVJLatHScvymf9eTYT6agZRuzQsFcnsqnFNiJaRNGzPufATpAaRwmKu88RH17qw2CLgjCjShMSV11HgnMzwTR4PaVlVQj0x5MqW+vVVEX2IEEqP+fL0lwo1",
}
var TestKeySet5 = []string{
	"ed25519 MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7",
}

func TestRSAPubKey(t *testing.T) {
	// Correct
	pk, err := PublicKeyFromString(TestKeySet1[0])
	assert.NoError(t, err)
	assert.Equal(t, "rsa", pk.Type.String())
	assert.Equal(t, "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAySB+eJGwb1Wu8KERTBLxKfMaZ7VRP92Q0LmJbqF5X2NrLyO+WkcYCcTqqoeW9unfoEiT4dS0I2YRLNcXztTaPs9xcoAWuXxPHwMKdjFKWQubttUoAAU08HpCkceO2y29NsrgOfKgUUc+D8FRXtyAe+GUg6wpWhxtT1BmjIgT4LIJKTKJ7cUdLxDQuKmT0uj0B+LQ4DTo4SvE7aV+Lb4wD2kwB0CHckHsbVC1jFPPcSLEtaSAjssaHe6v2c0oS2VdvuHq2pMa2lnuR/+wXeWn+zf0Rft0rcZNbIDpRt+DQH5g6VhA1Ww802mB/XYGaZcljzR5ArCDX8flvHP04ZcVWQIDAQAB", pk.S)
	assert.Equal(t, 65537, pk.K.(*rsa.PublicKey).E)
	assert.Equal(t, "", pk.Description)

	expectedErr := "incorrect key format"

	// Without type
	_, err = PublicKeyFromString(TestKeySet1[1])
	assert.EqualError(t, err, expectedErr)

	// Incorrect too
	_, err = PublicKeyFromString(TestKeySet1[2])
	assert.EqualError(t, err, expectedErr)

	// right type, wrong data
	_, err = PublicKeyFromString(TestKeySet1[3])
	assert.EqualError(t, err, expectedErr)

	// wrong type, right data
	_, err = PublicKeyFromString(TestKeySet1[4])
	assert.EqualError(t, err, "unsupported key type")

	pk, _ = PublicKeyFromString(TestKeySet1[0])
	npk, err := PublicKeyFromInterface(pk.Type, pk.K)
	assert.NoError(t, err)
	assert.Equal(t, "rsa", npk.Type.String())
	assert.Equal(t, "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAySB+eJGwb1Wu8KERTBLxKfMaZ7VRP92Q0LmJbqF5X2NrLyO+WkcYCcTqqoeW9unfoEiT4dS0I2YRLNcXztTaPs9xcoAWuXxPHwMKdjFKWQubttUoAAU08HpCkceO2y29NsrgOfKgUUc+D8FRXtyAe+GUg6wpWhxtT1BmjIgT4LIJKTKJ7cUdLxDQuKmT0uj0B+LQ4DTo4SvE7aV+Lb4wD2kwB0CHckHsbVC1jFPPcSLEtaSAjssaHe6v2c0oS2VdvuHq2pMa2lnuR/+wXeWn+zf0Rft0rcZNbIDpRt+DQH5g6VhA1Ww802mB/XYGaZcljzR5ArCDX8flvHP04ZcVWQIDAQAB", npk.S)
	assert.Equal(t, pk.K, npk.K.(*rsa.PublicKey))
	assert.Equal(t, "", npk.Description)
}

func TestED25519PubKey(t *testing.T) {
	key := ed25519.PublicKey{0x39, 0x3e, 0x4a, 0xf5, 0x5a, 0x91, 0x16, 0x32, 0x76, 0x43, 0x44, 0x56, 0x5a, 0xa4, 0xfc, 0xf8, 0xe9, 0xdd, 0x29, 0xa, 0xe6, 0xcd, 0x70, 0xae, 0x31, 0x34, 0xa7, 0x85, 0xfb, 0xd5, 0xa7, 0xb}
	// Correct ED25519
	pk, err := PublicKeyFromString(TestKeySet2[0])
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", pk.Type.String())
	assert.Equal(t, "MCowBQYDK2VwAyEAOT5K9VqRFjJ2Q0RWWqT8+OndKQrmzXCuMTSnhfvVpws=", pk.S)
	assert.Equal(t, key, pk.K.(ed25519.PublicKey))
	assert.Equal(t, "my foo bar", pk.Description)

	pk, err = PublicKeyFromInterface(pk.Type, key)
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", pk.Type.String())
	assert.Equal(t, "MCowBQYDK2VwAyEAOT5K9VqRFjJ2Q0RWWqT8+OndKQrmzXCuMTSnhfvVpws=", pk.S)
	assert.Equal(t, key, pk.K.(ed25519.PublicKey))
	assert.Equal(t, "", pk.Description)
}

func TestECDSAPub(t *testing.T) {
	x := new(big.Int)
	x.SetString("35674072579805598781611597844361315596081721500728226840731099278229644034819246176780429354491654135480285061184629", 10)
	y := new(big.Int)
	y.SetString("16152110512098221797044896547068689859301830683829073190715770376238845227548146305554820081396880025274135358569548", 10)

	pk, err := PublicKeyFromString(TestKeySet3[0])
	assert.NoError(t, err)
	assert.Equal(t, "ecdsa", pk.Type.String())
	assert.Equal(t, "MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAE58d01mIg3iWGqBHY7N6ch4L4LWya8Es2luWC08Wjn994nLgIOUp+cdMUfDYBe/x1aPE/yghGe3rrF4jW8uxVWy40BZK4NIu5yjMgSw0WBGTxOmZsVaA/xaOzvZSMTXxM", pk.S)
	assert.Equal(t, x, pk.K.(*ecdsa.PublicKey).X)
	assert.Equal(t, y, pk.K.(*ecdsa.PublicKey).Y)
	assert.Equal(t, "P-384", pk.K.(*ecdsa.PublicKey).Curve.Params().Name)
	assert.Equal(t, "", pk.Description)

	// Check public key data, but it's a private key
	_, err = PublicKeyFromString(TestKeySet4[0])
	assert.EqualError(t, err, "incorrect key format")
}

func TestECDSAPriv(t *testing.T) {
	d := new(big.Int)
	d.SetString("36622354286273408401742860803717975834272203860876773101980455299082624468593645020123023501111501921279897088878123", 10)

	pk, err := PrivateKeyFromString(TestKeySet4[0])
	assert.NoError(t, err)
	assert.Equal(t, "ecdsa", pk.Type.String())
	assert.Equal(t, "MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDDt8LNhN1AHVrBuBqKrryGgUS0EdRdjRuYGjONDYI96Vy8IcGsz4HQrX0biOzfE5iuhZANiAAT67cjQyt3qJUktq0dJy/KZ/15NhPpqBlG7NCwVyeyqcU2IlpE0bM+58BOkBpHCYq7zxEfXurDYIuCMKNKExJXXUeCczPBNHg9pWVVCPTHkypb69VURfYgQSo/58vSXCjU=", pk.S)
	assert.Equal(t, d, pk.K.(*ecdsa.PrivateKey).D)
	assert.Equal(t, "P-384", pk.K.(*ecdsa.PrivateKey).Curve.Params().Name)
}

func TestED25519Priv(t *testing.T) {
	b := []byte{0xc7, 0x54, 0xb9, 0x99, 0x3c, 0x19, 0x73, 0xc, 0x46, 0x7a, 0xdb, 0x67, 0xdd, 0x9, 0xc, 0xd2, 0x2b, 0xcf, 0x98, 0x7b, 0x69, 0xcd, 0xb0, 0x85, 0x98, 0xaa, 0xa6, 0xf7, 0x83, 0x4a, 0x6, 0x7b}

	pk, err := PrivateKeyFromString(TestKeySet5[0])
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", pk.Type.String())
	assert.Equal(t, "MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7", pk.S)
	assert.Equal(t, b, pk.K.(ed25519.PrivateKey).Seed())

	npk, err := PrivateKeyFromInterface(pk.Type, pk.K)
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", pk.Type.String())
	assert.Equal(t, "MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7", pk.S)
	assert.Equal(t, pk.K, npk.K.(ed25519.PrivateKey))
}

func TestIncorrectKeys(t *testing.T) {
	pk, err := PrivateKeyFromString("foo 12314")
	assert.Error(t, err)
	assert.Nil(t, pk)

	pk, err = PrivateKeyFromString("foo12314")
	assert.Error(t, err)
	assert.Nil(t, pk)
}

func TestPrivateKeyJSON(t *testing.T) {
	pk, err := PrivateKeyFromString(TestKeySet5[0])
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", pk.Type.String())

	data, err := pk.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "\"ed25519 MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7\"", string(data))

	err = pk.UnmarshalJSON([]byte("fasfdsadfA"))
	assert.IsType(t, &json.SyntaxError{}, err)

	foo := &PrivKey{}
	err = foo.UnmarshalJSON([]byte("\"ed25519 MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7\""))
	assert.NoError(t, err)
	assert.Equal(t, "ed25519", foo.Type.String())
	assert.Equal(t, "MC4CAQAwBQYDK2VwBCIEIMdUuZk8GXMMRnrbZ90JDNIrz5h7ac2whZiqpveDSgZ7", foo.S)
}

func TestPubKeyJSON(t *testing.T) {
	// Pub keys
	pk, err := PublicKeyFromString(TestKeySet1[0])
	assert.NoError(t, err)
	assert.Equal(t, "rsa", pk.Type.String())

	data, err := pk.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "\"rsa MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAySB+eJGwb1Wu8KERTBLxKfMaZ7VRP92Q0LmJbqF5X2NrLyO+WkcYCcTqqoeW9unfoEiT4dS0I2YRLNcXztTaPs9xcoAWuXxPHwMKdjFKWQubttUoAAU08HpCkceO2y29NsrgOfKgUUc+D8FRXtyAe+GUg6wpWhxtT1BmjIgT4LIJKTKJ7cUdLxDQuKmT0uj0B+LQ4DTo4SvE7aV+Lb4wD2kwB0CHckHsbVC1jFPPcSLEtaSAjssaHe6v2c0oS2VdvuHq2pMa2lnuR/+wXeWn+zf0Rft0rcZNbIDpRt+DQH5g6VhA1Ww802mB/XYGaZcljzR5ArCDX8flvHP04ZcVWQIDAQAB\"", string(data))

	err = pk.UnmarshalJSON([]byte("fasfdsadfA"))
	assert.IsType(t, &json.SyntaxError{}, err)

	foo := &PubKey{}
	err = foo.UnmarshalJSON([]byte("\"rsa MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8pvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5aG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxXPOhC1QUcllt08yBB3wIDAQAB\""))
	assert.NoError(t, err)
	assert.Equal(t, "rsa", foo.Type.String())
	assert.Equal(t, "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC57qC/BeoYcM6ijazuaCdJkbT8pvPpFEDVzf9ZQ9axswXU3mywSOaR3wflriSjmvRfUNs/BAjshgtJqgviUXx7lE5aG9mcUyvomyFFpfCR2l2Lvow0H8y7JoL6yxMSQf8gpAcaQzPB8dsfGe+DqA+5wjxXPOhC1QUcllt08yBB3wIDAQAB", foo.S)
}

func TestFingerprint(t *testing.T) {
	pk, _ := PublicKeyFromString(TestKeySet1[0])
	assert.Equal(t, "4d6dc360edab1276404d3ece66868d6b390288b5e4bad42c75b350b792127046", pk.Fingerprint())

	pk, _ = PublicKeyFromString(TestKeySet2[0])
	assert.Equal(t, "3bde10ad2dc6508163a8a4944c46684102ed777e57f73514c75e6e62448a3d85", pk.Fingerprint())
}
