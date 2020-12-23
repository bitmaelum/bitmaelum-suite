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
	"fmt"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Authenticate holds multiple middleware/authenticators that can authenticate against the API
type Authenticate struct {
	Chain []Authenticator
}

// ContextKeyType is the type used in context
type ContextKeyType int

// Context keys
const (
	AuthorizationContext ContextKeyType = iota // Defines which authorization type is used
	APIKeyContext                              // Used API key
	AuthKeyContext                             // Used Auth key
	ClaimsContext                              // Context key for fetching JWT claims
	AddressContext                             // Context key for fetching the address from the JWT
)

// AuthStatus is the status that is returned by an authenticator
type AuthStatus int

// AuthorizationContext is a context key that defines the authorization method
const (
	AuthStatusPass    AuthStatus = iota // AuthStatusPass is returned when the authenticator cannot handle or doesn't care. Next authenticator is checked
	AuthStatusSuccess                   // AuthStatusSuccess allows the request, no other authenticators are checked
	AuthStatusFailure                   // AuthStatusFailure denies the request, no other authenticators are checked
)

// Authenticator allows you to use the struct in the multi-auth middleware
type Authenticator interface {
	Authenticate(req *http.Request, route string) (AuthStatus, context.Context, error)
}

// Middleware JWT token authentication
func (ma *Authenticate) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// We need to fetch the current route here, because if we fetch it from too deep, we get into trouble with
		// testing the given components.
		currentRouteName := mux.CurrentRoute(req).GetName()

		for _, auth := range ma.Chain {
			logrus.Tracef("authenticate: trying %T", auth)
			status, ctx, err := auth.Authenticate(req, currentRouteName)

			if status == AuthStatusPass {
				// Continue with next authenticator
				continue
			}

			if status == AuthStatusFailure {
				logrus.Tracef("authenticate: failure: %s", err)
				ErrorOut(w, http.StatusUnauthorized, err.Error())
				return
			}

			if status == AuthStatusSuccess {
				ctx = context.WithValue(ctx, AuthorizationContext, fmt.Sprintf("%T", auth))
				logrus.Tracef("authenticate: found ok %T", auth)
				next.ServeHTTP(w, req.WithContext(ctx))
				return
			}

			return
		}

		logrus.Tracef("authenticate: unauthorized")
		ErrorOut(w, http.StatusUnauthorized, "Unauthorized")
	})
}

// Add a new authenticator to the list
func (ma *Authenticate) Add(auth Authenticator) {
	ma.Chain = append(ma.Chain, auth)
}

// ErrorOut outputs an error
func ErrorOut(w http.ResponseWriter, code int, msg string) {
	type OutputResponse struct {
		Error  bool   `json:"error,omitempty"`
		Status string `json:"status"`
	}

	logrus.Debugf("Returning error (%d): %s", code, msg)

	_ = httputils.JSONOut(w, code, &OutputResponse{
		Error:  true,
		Status: msg,
	})
}
