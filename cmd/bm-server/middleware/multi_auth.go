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
