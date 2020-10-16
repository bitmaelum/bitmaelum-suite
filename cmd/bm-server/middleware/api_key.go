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

package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/apikey"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
)

// APIKey is a middleware that automatically verifies given API key
type APIKey struct{}

// ErrInvalidAPIKey is returned when the API key is not valid
var ErrInvalidAPIKey = errors.New("invalid API key")

// ErrExpiredAPIKey is returned when a key is expired
var ErrExpiredAPIKey = errors.New("expired API key")

// ErrInvalidAuthentication is returned when no or invalid authentication method is found
var ErrInvalidAuthentication = errors.New("invalid authentication")

// APIKeyContext is a context key with the value the API key
type APIKeyContext string

// Middleware API token authentication
func (mw *APIKey) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx, ok := mw.Authenticate(req)
		if !ok {
			ErrorOut(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

// @TODO make sure we can't use a key to fetch other people's info

// Authenticate will check if an API key matches the request
func (*APIKey) Authenticate(req *http.Request) (context.Context, bool) {
	key, err := getAPIKey(req.Header.Get("Authorization"))
	if err != nil {
		return nil, false
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, APIKeyContext("api-key"), key)

	return ctx, true
}

// check authorization API key
func getAPIKey(auth string) (*apikey.KeyType, error) {
	if auth == "" {
		return nil, ErrInvalidAuthentication
	}

	if len(auth) <= 6 || strings.ToUpper(auth[0:7]) != "BEARER " {
		return nil, ErrInvalidAuthentication
	}
	apiKeyID := auth[7:]

	ar := container.GetAPIKeyRepo()
	key, err := ar.Fetch(apiKeyID)
	if err != nil {
		return nil, ErrInvalidAPIKey
	}

	if !key.ValidUntil.IsZero() && time.Now().After(key.ValidUntil) {
		return nil, ErrExpiredAPIKey
	}

	return key, nil
}
