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

package handler

import (
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
)

// GetAPIKey returns the api key stored in the request context. returns nil when not found
func GetAPIKey(req *http.Request) *key.APIKeyType {
	return req.Context().Value(middleware.APIKeyContext).(*key.APIKeyType)
}

// IsAPIKeyAuthenticated returns true when the given request is authenticated by a api key
func IsAPIKeyAuthenticated(req *http.Request) bool {
	return req.Context().Value(middleware.AuthorizationContext) == "*auth.APIKeyAuth"
}

// GetAuthKey will return the auth key if any. Returns nil when not found
func GetAuthKey(req *http.Request) *key.AuthKeyType {
	return req.Context().Value(middleware.AuthKeyContext).(*key.AuthKeyType)
}

// IsAuthKeyAuthenticated returns true when the given request is authenticated by a auth key
func IsAuthKeyAuthenticated(req *http.Request) bool {
	return req.Context().Value(middleware.AuthorizationContext) == "*auth.OnBehalfJwtAuth"
}

