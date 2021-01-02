// Copyright (c) 2021 BitMaelum Authors
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
	"crypto/subtle"
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/sirupsen/logrus"
	"github.com/vtolstov/jwt-go"
)

// Error codes
var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrTokenNotValid        = errors.New("token not valid")
	ErrTokenTimeNotValid    = errors.New("token time not valid")
	ErrSubjectNotValid      = errors.New("subject not valid")
)

/*
 * @TODO
 * I don't really like this. Suppose we get access to a single JWT token. We can use the same token for every
 * single call in the next hour. Maybe we should limit each token for single-use (with a nonce that expires after
 * one hour, which means that expiresAt expires too), or maybe even limit the token for a single operation (add
 * request info?)
 */

const (
	invalidSigningMethod string = "invalid signing method"
)

// GenerateJWTToken generates a JWT token with the address and singed by the given private key
func GenerateJWTToken(addr hash.Hash, key bmcrypto.PrivKey) (string, error) {
	// Get the current 30 second block of time
	now := jwt.TimeFunc()
	ct := (now.Unix() / 30) * 30

	claims := &jwt.StandardClaims{
		ExpiresAt: ct + 30 + 30, // Expires after the NEXT 30 second block
		NotBefore: ct - 30,      // Start accepting in the previous 30 second block
		IssuedAt:  ct - 30,      // IssuesAt cannot be more than NotBefore
		Subject:   addr.String(),
	}

	token := jwt.NewWithClaims(key.Type.JWTSignMethod(), claims)

	return token.SignedString(key.K)
}

// IsJWTTokenExpired will check if the token is already expired
func IsJWTTokenExpired(tokenString string) bool {
	// Get Claims from Token
	claims := &jwt.StandardClaims{}
	new(jwt.Parser).ParseUnverified(tokenString, claims)

	// Calculate current time sliding window
	now := jwt.TimeFunc()
	ct := (now.Unix() / 30) * 30

	return claims.IssuedAt != ct-30
}

// ValidateJWTToken validates a JWT token with the given public key and address
func ValidateJWTToken(tokenString string, addr hash.Hash, key bmcrypto.PubKey) (*jwt.Token, error) {
	logrus.Tracef("validating JWT token")

	// Just return the key from the token
	kf := func(token *jwt.Token) (interface{}, error) {
		return key.K, nil
	}

	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, kf)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			// When timing is off, return a separate message, as we can pass this to the caller
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				logrus.Trace("auth: jwt: token time is not valid")
				return nil, ErrTokenTimeNotValid
			}
		}

		logrus.Trace("auth: jwt: ", err)
		return nil, err
	}

	// Make sure the signature method of the JWT matches our public key
	if !key.Type.JWTHasValidSignMethod(token) {
		logrus.Tracef("auth: jwt: " + invalidSigningMethod)
		return nil, ErrInvalidSigningMethod
	}

	// It should be a valid token
	if !token.Valid {
		logrus.Trace("auth: jwt: token not valid")
		return nil, ErrTokenNotValid
	}

	// The standard claims should be valid
	err = token.Claims.Valid()
	if err != nil {
		logrus.Trace("auth: jwt: ", err)
		return nil, err
	}

	// Check subject explicitly
	res := subtle.ConstantTimeCompare([]byte(token.Claims.(*jwt.StandardClaims).Subject), []byte(addr.String()))
	if res == 0 {
		logrus.Tracef("auth: jwt: subject does not match")
		return nil, ErrSubjectNotValid
	}

	logrus.Trace("auth: jwt: token is valid")
	return token, nil
}

func init() {
	var edDSASigningMethod bmcrypto.SigningMethodEdDSA
	jwt.RegisterSigningMethod(edDSASigningMethod.Alg(), func() jwt.SigningMethod { return &edDSASigningMethod })
}
