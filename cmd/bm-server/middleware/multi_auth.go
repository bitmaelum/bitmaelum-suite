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

	"github.com/sirupsen/logrus"
)

// MultiAuth holds multiple middleware/authenticators that can authenticate against the API
type MultiAuth struct {
	Auths []Authenticable
}

type authContext string

// Authenticable allows you to use the struct in the multi-auth middleware
type Authenticable interface {
	Authenticate(req *http.Request) (context.Context, bool)
}

// Middleware JWT token authentication
func (ma *MultiAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, auth := range ma.Auths {
			logrus.Tracef("multiauth. Trying %T", auth)
			ctx, ok := auth.Authenticate(req)
			if ok {
				ctx = context.WithValue(ctx, authContext("auth_method"), fmt.Sprintf("%T", auth))
				logrus.Tracef("multiauth found ok %T", auth)
				next.ServeHTTP(w, req.WithContext(ctx))
				return
			}
		}

		logrus.Tracef("multiauth unauthorized")
		ErrorOut(w, http.StatusUnauthorized, "Unauthorized")
	})
}
