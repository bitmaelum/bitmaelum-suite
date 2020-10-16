// Copyright (c) 2020 BitMaelum Authors
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

package internal

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/vtolstov/jwt-go"
)

const (
	mockToken     = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODU2OTYsImlhdCI6MTU3Nzg4MjA5NiwibmJmIjoxNTc3ODgyMDk2LCJzdWIiOiIxODgyYjkxYjdmNDlkNDc5Y2YxZWMyZjFlY2VlMzBkMGU1MzkyZTk2M2EyMTA5MDE1YjcxNDliZjcxMmFkMWI2In0.e9osJ5LRAkz6hgFMb9hSJe9SDcefi3l3t7q5NGXm4BisNKQa0lfmefVZXAdi5U8PKD3laSFjtkgIoN97TWc5o7b4KxFPziAb1KZ0JHz0oD8MBjFf0ebrlWv5GLsEozyFfyID9onvOVI6purY4ZBmiap3ncp2gHip7KFZQweVcR4"
	mockSignature = "e9osJ5LRAkz6hgFMb9hSJe9SDcefi3l3t7q5NGXm4BisNKQa0lfmefVZXAdi5U8PKD3laSFjtkgIoN97TWc5o7b4KxFPziAb1KZ0JHz0oD8MBjFf0ebrlWv5GLsEozyFfyID9onvOVI6purY4ZBmiap3ncp2gHip7KFZQweVcR4"
)

func TestGenerateJWTToken(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/privkey.rsa")
	privKey, err := bmcrypto.NewPrivKey(string(data))
	assert.Nil(t, err)

	haddr := hash.New("test!")

	token, err := GenerateJWTToken(haddr, *privKey)
	assert.Nil(t, err)
	assert.Equal(t, mockToken, token)
}

func TestValidateJWTToken(t *testing.T) {
	data, _ := ioutil.ReadFile("../testdata/pubkey.rsa")
	pubKey, _ := bmcrypto.NewPubKey(string(data))

	haddr := hash.New("test!")

	token, err := ValidateJWTToken(mockToken, haddr, *pubKey)
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
