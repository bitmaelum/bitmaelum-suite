package internal

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/stretchr/testify/assert"
	"github.com/vtolstov/jwt-go"
)

const (
	mockToken     = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODU2OTYsImlhdCI6MTU3Nzg4MjA5NiwibmJmIjoxNTc3ODgyMDk2LCJzdWIiOiJhNjlhNDM1NjM4MjViOTU3MzRlYWQ1NTVmNDg1YjdiYTM5M2QyOTEzNTYwMjUwOWQxYWRkYWQ3N2QyNGYxMDUxIn0.jtnAY2WUfQHBtDGduWQaig25c1uQClYXPEkYoU5cXkSQaiewoR1sU9zRLctEPN1nKuTSig6SNnXPrkBGEOg3Z69WlldubklG8k_f5DSZ3qjiWxS_mDiGhAqWhjBMe-IBWvp8oiblEqV2upRfR89XcMKHbBEQ20awrdSbI5zXbFw"
	mockSignature = "jtnAY2WUfQHBtDGduWQaig25c1uQClYXPEkYoU5cXkSQaiewoR1sU9zRLctEPN1nKuTSig6SNnXPrkBGEOg3Z69WlldubklG8k_f5DSZ3qjiWxS_mDiGhAqWhjBMe-IBWvp8oiblEqV2upRfR89XcMKHbBEQ20awrdSbI5zXbFw"
)

func TestGenerateJWTToken(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/privkey.rsa")
	privKey, err := bmcrypto.NewPrivKey(string(data))
	assert.Nil(t, err)

	haddr, _ := address.NewHash("test!")

	token, err := GenerateJWTToken(*haddr, *privKey)
	assert.Nil(t, err)
	assert.Equal(t, mockToken, token)
}

func TestValidateJWTToken(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/pubkey.rsa")
	pubKey, _ := bmcrypto.NewPubKey(string(data))

	haddr, _ := address.NewHash("test!")

	token, err := ValidateJWTToken(mockToken, *haddr, *pubKey)
	assert.Nil(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, "RS256", token.Method.Alg())
	assert.Equal(t, mockSignature, token.Signature)
	assert.Equal(t, int64(1577882096), token.Claims.(*jwt.StandardClaims).IssuedAt)
	assert.Equal(t, int64(1577885696), token.Claims.(*jwt.StandardClaims).ExpiresAt)
	assert.Equal(t, int64(1577882096), token.Claims.(*jwt.StandardClaims).NotBefore)
	assert.Equal(t, haddr.String(), token.Claims.(*jwt.StandardClaims).Subject)
}

func init() {
	// Mock JWT time
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 56, 0, time.UTC)
	}
}
