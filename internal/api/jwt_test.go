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

package api

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
	mockToken     = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc4ODIxMzAsImlhdCI6MTU3Nzg4MjA0MCwibmJmIjoxNTc3ODgyMDQwLCJzdWIiOiIxODgyYjkxYjdmNDlkNDc5Y2YxZWMyZjFlY2VlMzBkMGU1MzkyZTk2M2EyMTA5MDE1YjcxNDliZjcxMmFkMWI2In0.D0QCl93sfwtOVmHpt5LJ9OjnLfNR0d9WnyZIVVa-Ktxd-PSLC6b-UlhSV3NKnMz1mNdO3KQIf9_0RQcjWWxOUzH2kXANMNngeLz5bHQowiSDtTMFwKdwCdHhMaYCuMkEGvILfKRUuDhussSFZmcGcqDkIqKRFDN-0HyoHfCEHo4"
	mockSignature = "D0QCl93sfwtOVmHpt5LJ9OjnLfNR0d9WnyZIVVa-Ktxd-PSLC6b-UlhSV3NKnMz1mNdO3KQIf9_0RQcjWWxOUzH2kXANMNngeLz5bHQowiSDtTMFwKdwCdHhMaYCuMkEGvILfKRUuDhussSFZmcGcqDkIqKRFDN-0HyoHfCEHo4"
)

func TestGenerateJWTToken(t *testing.T) {
	// Mock JWT time
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	data, _ := ioutil.ReadFile("../../testdata/privkey.rsa")
	privKey, err := bmcrypto.PrivateKeyFromString(string(data))
	assert.Nil(t, err)

	haddr := hash.New("test!")

	token, err := GenerateJWTToken(haddr, *privKey)
	assert.Nil(t, err)
	assert.Equal(t, mockToken, token)
}

func TestValidateJwtTokenExpiry(t *testing.T) {
	data, _ := ioutil.ReadFile("../../testdata/pubkey.rsa")
	pubKey, _ := bmcrypto.PublicKeyFromString(string(data))
	haddr := hash.New("test!")

	// The current time block in our mock token is 12:34:30 - 12:36:00 (1577882040 - 1577882130)

	// In previous 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 29, 0, time.UTC)
	}
	token, err := ValidateJWTToken(mockToken, haddr, *pubKey)
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)

	// Before previous 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 33, 59, 0, time.UTC)
	}
	token, err = ValidateJWTToken(mockToken, haddr, *pubKey)
	assert.Error(t, err)
	assert.Nil(t, token)

	// In current 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 59, 0, time.UTC)
	}
	token, err = ValidateJWTToken(mockToken, haddr, *pubKey)
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)

	// In next 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 35, 29, 0, time.UTC)
	}
	token, err = ValidateJWTToken(mockToken, haddr, *pubKey)
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.True(t, token.Valid)

	// After next 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 35, 31, 0, time.UTC)
	}
	token, err = ValidateJWTToken(mockToken, haddr, *pubKey)
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestValidateJWTToken(t *testing.T) {
	// Mock JWT time
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 56, 0, time.UTC)
	}

	data, _ := ioutil.ReadFile("../../testdata/pubkey.rsa")
	pubKey, _ := bmcrypto.PublicKeyFromString(string(data))

	haddr := hash.New("test!")

	token, err := ValidateJWTToken(mockToken, haddr, *pubKey)
	assert.Nil(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, "RS256", token.Method.Alg())
	assert.Equal(t, mockSignature, token.Signature)
	assert.Equal(t, int64(1577882040), token.Claims.(*jwt.StandardClaims).IssuedAt)
	assert.Equal(t, int64(1577882130), token.Claims.(*jwt.StandardClaims).ExpiresAt)
	assert.Equal(t, int64(1577882040), token.Claims.(*jwt.StandardClaims).NotBefore)
	assert.Equal(t, haddr.String(), token.Claims.(*jwt.StandardClaims).Subject)
}

func TestIsJWTTokenExpired(t *testing.T) {
	// The current time block in our mock token is 12:34:30 - 12:36:00 (1577882040 - 1577882130)
	// This test should only be valid if it is on current time block

	// In previous 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 29, 0, time.UTC)
	}
	assert.True(t, IsJWTTokenExpired(mockToken))

	// Before previous 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 33, 59, 0, time.UTC)
	}
	assert.True(t, IsJWTTokenExpired(mockToken))

	// In current 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 34, 59, 0, time.UTC)
	}
	assert.False(t, IsJWTTokenExpired(mockToken))

	// In next 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 35, 29, 0, time.UTC)
	}
	assert.True(t, IsJWTTokenExpired(mockToken))

	// After next 30 second block
	jwt.TimeFunc = func() time.Time {
		return time.Date(2020, 01, 01, 12, 35, 31, 0, time.UTC)
	}
	assert.True(t, IsJWTTokenExpired(mockToken))
}

func TestED25519(t *testing.T) {
	kt, err := bmcrypto.FindKeyType("ed25519")
	assert.NoError(t, err)
	priv, pub, err := bmcrypto.GenerateKeyPair(kt)
	assert.NoError(t, err)

	tokenStr, err := GenerateJWTToken(hash.New("test!"), *priv)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	token, err := ValidateJWTToken(tokenStr, hash.New("test!"), *pub)
	assert.NoError(t, err)
	assert.Equal(t, "EdDSA", token.Method.Alg())
	assert.Equal(t, "1882b91b7f49d479cf1ec2f1ecee30d0e5392e963a2109015b7149bf712ad1b6", token.Claims.(*jwt.StandardClaims).Subject)
}

func TestECDSA(t *testing.T) {
	kt, err := bmcrypto.FindKeyType("ecdsa")
	assert.NoError(t, err)
	priv, pub, err := bmcrypto.GenerateKeyPair(kt)
	assert.NoError(t, err)

	tokenStr, err := GenerateJWTToken(hash.New("test!"), *priv)
	assert.NoError(t, err)

	token, err := ValidateJWTToken(tokenStr, hash.New("test!"), *pub)
	assert.NoError(t, err)
	assert.Equal(t, "ES384", token.Method.Alg())
	assert.Equal(t, "1882b91b7f49d479cf1ec2f1ecee30d0e5392e963a2109015b7149bf712ad1b6", token.Claims.(*jwt.StandardClaims).Subject)
}

func init() {
	var edDSASigningMethod bmcrypto.SigningMethodEdDSA
	jwt.RegisterSigningMethod(edDSASigningMethod.Alg(), func() jwt.SigningMethod { return &edDSASigningMethod })
}
