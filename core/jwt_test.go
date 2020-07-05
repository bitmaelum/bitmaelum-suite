package core

import (
	"github.com/bitmaelum/bitmaelum-server/core/encrypt"
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
	"github.com/stretchr/testify/assert"
	"github.com/vtolstov/jwt-go"
	"io/ioutil"
	"testing"
	"time"
)

const (
	mockToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODU2OTYsImlhdCI6MTU3Nzg4MjA5NiwibmJmIjoxNTc3ODgyMDk2LCJzdWIiOiIxODgyYjkxYjdmNDlkNDc5Y2YxZWMyZjFlY2VlMzBkMGU1MzkyZTk2M2EyMTA5MDE1YjcxNDliZjcxMmFkMWI2In0.e9osJ5LRAkz6hgFMb9hSJe9SDcefi3l3t7q5NGXm4BisNKQa0lfmefVZXAdi5U8PKD3laSFjtkgIoN97TWc5o7b4KxFPziAb1KZ0JHz0oD8MBjFf0ebrlWv5GLsEozyFfyID9onvOVI6purY4ZBmiap3ncp2gHip7KFZQweVcR4"
	mockSignature = "e9osJ5LRAkz6hgFMb9hSJe9SDcefi3l3t7q5NGXm4BisNKQa0lfmefVZXAdi5U8PKD3laSFjtkgIoN97TWc5o7b4KxFPziAb1KZ0JHz0oD8MBjFf0ebrlWv5GLsEozyFfyID9onvOVI6purY4ZBmiap3ncp2gHip7KFZQweVcR4"
)



func TestGenerateJWTToken(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/mykey.pem")
	privKey, _ := encrypt.PEMToPrivKey(data)

	haddr, _ := address.NewHash("test!")

	token, err := GenerateJWTToken(*haddr, privKey)
	assert.Nil(t, err)
	assert.Equal(t, mockToken, token)
}

func TestValidateJWTToken(t *testing.T) {
	data, _ := ioutil.ReadFile("./testdata/mykey.pub")
	pubKey, _ := encrypt.PEMToPubKey(data)

	haddr, _ := address.NewHash("test!")

	token, err := ValidateJWTToken(mockToken, *haddr, pubKey)
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
