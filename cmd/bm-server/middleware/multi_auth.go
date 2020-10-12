package middleware

import (
	"context"
	"net/http"
)

// MultiAuth holds multiple middleware/authenticators that can authenticate against the API
type MultiAuth struct {
	Auths []Authenticable
}

// Authenticable allows you to use the struct in the multi-auth middleware
type Authenticable interface {
	Authenticate(req *http.Request) (context.Context, bool)
}

// Middleware JWT token authentication
func (ma *MultiAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, auth := range ma.Auths {
			ctx, ok := auth.Authenticate(req)
			if ok {
				next.ServeHTTP(w, req.WithContext(ctx))
			}
		}

		ErrorOut(w, http.StatusUnauthorized, "Unauthorized")
	})
}
