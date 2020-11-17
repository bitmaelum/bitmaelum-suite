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
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIKeyAuth is a middleware that automatically verifies given API key
type APIKeyAuth struct {
	PermissionList map[string][]string
}

var (
	// ErrInvalidAPIKey is returned when the API key is not valid
	ErrInvalidAPIKey = errors.New("invalid API key")
	// ErrExpiredAPIKey is returned when a key is expired
	ErrExpiredAPIKey = errors.New("expired API key")
	// ErrInvalidAuthentication is returned when no or invalid authentication method is found
	ErrInvalidAuthentication = errors.New("invalid authentication")
	// ErrIncorrectRoute is returned when a route without a name is used (permissions are checked against named routes only)
	ErrIncorrectRoute = errors.New("api keys need named routes")
	// ErrInvalidPermission When a user does not have the correct permission
	ErrInvalidPermission = errors.New("api keys need named routes")
)

type contextKey int

const (
	// APIKeyContext is a context key with the value the API key
	APIKeyContext contextKey = iota
)

// @TODO make sure we can't use a key to fetch other people's info

// Authenticate will check if an API key matches the request
func (a *APIKeyAuth) Authenticate(req *http.Request, route string) (context.Context, bool) {
	// Check if the address actually exists
	haddr, err := hash.NewFromHash(mux.Vars(req)["addr"])
	if err != nil {
		return nil, false
	}

	accountRepo := container.Instance.GetAccountRepo()
	if !accountRepo.Exists(*haddr) {
		logrus.Trace("auth: address not found")
		return nil, false
	}

	// Check api key.
	k, err := a.checkAPIKey(req.Header.Get("Authorization"), *haddr, route)
	if err != nil {
		return nil, false
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, APIKeyContext, k)

	return ctx, true
}

func (a *APIKeyAuth) checkAPIKey(bearerToken string, addrHash hash.Hash, routeName string) (*key.APIKeyType, error) {
	k, err := a.getAPIKey(bearerToken)
	if err != nil {
		return nil, err
	}

	if k.AddressHash.String() != addrHash.String() {
		return nil, ErrInvalidAPIKey
	}

	if !k.Expires.IsZero() && time.Now().After(k.Expires) {
		return nil, ErrExpiredAPIKey
	}

	// Check permissions
	if routeName == "" {
		return nil, ErrIncorrectRoute
	}

	perms, ok := a.PermissionList[routeName]
	if !ok {
		return nil, ErrInvalidPermission
	}

	for _, perm := range perms {
		for _, userperm := range k.Permissions {
			if userperm == perm {
				return k, nil
			}
		}
	}

	return nil, ErrInvalidPermission
}

// check authorization API key
func (*APIKeyAuth) getAPIKey(bearerToken string) (*key.APIKeyType, error) {
	if bearerToken == "" {
		return nil, ErrInvalidAuthentication
	}

	if len(bearerToken) <= 6 || strings.ToUpper(bearerToken[0:7]) != "BEARER " {
		return nil, ErrInvalidAuthentication
	}
	apiKeyID := bearerToken[7:]

	apiKeyRepo := container.Instance.GetAPIKeyRepo()
	k, err := apiKeyRepo.Fetch(apiKeyID)
	if err != nil {
		return nil, ErrInvalidAPIKey
	}

	return k, nil
}
