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

// Authenticate will check if an API key matches the request
func (*APIKey) Authenticate(req *http.Request) (context.Context, bool) {
	key, err := getAPIKey(req.Header.Get("Authorization"))
	if err != nil {
		return nil, false
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, APIKeyContext("apikey"), key)

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
