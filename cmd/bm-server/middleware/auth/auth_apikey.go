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
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIKeyAuth is a middleware that automatically verifies given API key
type APIKeyAuth struct {
	PermissionList map[string][]string // Map of route -> permission mappings
	AdminKeys      bool                // True when the authorizer checks admin keys instead of account keys
}

var (
	errInvalidAPIKey         = errors.New("invalid API key")
	errExpiredAPIKey         = errors.New("expired API key")
	errInvalidAuthentication = errors.New("invalid authentication")
	errIncorrectRoute        = errors.New("api keys need named routes")
	errInvalidPermission     = errors.New("api keys need named routes")
)

type contextKey int

const (
	// APIKeyContext is a context key with the value the API key
	APIKeyContext contextKey = iota
)

// Authenticate will check if an API key matches the request
func (a *APIKeyAuth) Authenticate(req *http.Request, route string) (middleware.AuthStatus, context.Context, error) {
	var haddr *hash.Hash = nil
	if !a.AdminKeys {
		// Check if the address actually exists
		var err error
		haddr, err = hash.NewFromHash(mux.Vars(req)["addr"])
		if err != nil {
			logrus.Trace("auth: addr not found in url")
			return middleware.AuthStatusPass, nil, nil
		}

		accountRepo := container.Instance.GetAccountRepo()
		if !accountRepo.Exists(*haddr) {
			logrus.Trace("auth: address not found")
			return middleware.AuthStatusPass, nil, nil
		}
	}

	// Check api key.
	k, err := a.checkAPIKey(req.Header.Get("Authorization"), haddr, route)
	if err != nil {
		logrus.Trace("auth: checkApiKey failed: ", err)
		return middleware.AuthStatusPass, nil, nil
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, APIKeyContext, k)

	logrus.Trace("auth: checkApiKey success: ", err)
	return middleware.AuthStatusSuccess, ctx, nil
}

func (a *APIKeyAuth) checkAPIKey(bearerToken string, addrHash *hash.Hash, routeName string) (*key.APIKeyType, error) {
	k, err := a.getAPIKey(bearerToken)
	if err != nil {
		logrus.Trace("auth: can't get api key: ", err)
		return nil, err
	}

	// Check if the address hash matches the hash of the key, but only if the hash is given
	if !a.AdminKeys && k.AddressHash.String() != addrHash.String() {
		return nil, errInvalidAPIKey
	}

	if !k.Expires.IsZero() && time.Now().After(k.Expires) {
		logrus.Trace("auth: checkApiKey expired key")
		return nil, errExpiredAPIKey
	}

	// Check permissions
	if routeName == "" {
		logrus.Trace("auth: checkApiKey no route found")
		return nil, errIncorrectRoute
	}

	logrus.Tracef("checking permission of route %s", routeName)
	perms, ok := a.PermissionList[routeName]
	if !ok {
		return nil, errInvalidPermission
	}

	for _, perm := range perms {
		for _, userperm := range k.Permissions {
			if userperm == perm {
				return k, nil
			}
		}
	}

	logrus.Trace("auth: checkApiKey no permission found")
	return nil, errInvalidPermission
}

// check authorization API key
func (*APIKeyAuth) getAPIKey(bearerToken string) (*key.APIKeyType, error) {
	if bearerToken == "" {
		return nil, errInvalidAuthentication
	}

	if len(bearerToken) <= 6 || strings.ToUpper(bearerToken[0:7]) != "BEARER " {
		return nil, errInvalidAuthentication
	}
	apiKeyID := bearerToken[7:]

	apiKeyRepo := container.Instance.GetAPIKeyRepo()
	k, err := apiKeyRepo.Fetch(apiKeyID)
	if err != nil {
		return nil, errInvalidAPIKey
	}

	return k, nil
}
