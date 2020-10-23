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

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/vtolstov/jwt-go"
)

// AuthJwt is a middleware that automatically verifies given JWT token
type AuthJwt struct{}

type ClaimsContext string
type AddressContext string

// ErrTokenNotValidated is returned when the token could not be validated (for any reason)
var ErrTokenNotValidated = errors.New("token could not be validated")


var accountRepo = container.GetAccountRepo()

// Authenticate will check if an API key matches the request
func (mw *AuthJwt) Authenticate(req *http.Request, _ string) (context.Context, bool) {
	// Check if the address actually exists
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		return nil, false
	}

	if !accountRepo.Exists(*haddr) {
		logrus.Trace("auth: address not found")
		return nil, false
	}

	// Check token
	token, err := checkToken(req.Header.Get("Authorization"), *haddr)
	if err != nil {
		logrus.Trace("auth: incorrect token: ", err)
		return nil, false
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, ClaimsContext("claims"), token.Claims)
	ctx = context.WithValue(ctx, AddressContext("address"), token.Claims.(*jwt.StandardClaims).Subject)

	return ctx, true
}

// Check if the authorization contains a valid JWT token for the given address
func checkToken(bearerToken string, addr hash.Hash) (*jwt.Token, error) {
	if bearerToken == "" {
		logrus.Trace("auth: empty auth string")
		return nil, ErrTokenNotValidated
	}

	if len(bearerToken) <= 6 || strings.ToUpper(bearerToken[0:7]) != "BEARER " {
		logrus.Trace("auth: bearer not found")
		return nil, ErrTokenNotValidated
	}
	tokenString := bearerToken[7:]

	keys, err := accountRepo.FetchKeys(addr)
	if err != nil {
		logrus.Trace("auth: cannot fetch keys: ", err)
		return nil, ErrTokenNotValidated
	}

	for _, key := range keys {
		token, err := internal.ValidateJWTToken(tokenString, addr, key)
		if err == nil {
			return token, nil
		}
	}

	logrus.Trace("auth: no key found that validates the token")
	return nil, ErrTokenNotValidated
}
