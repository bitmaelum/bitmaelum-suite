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
func (*APIKey) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key, err := getAPIKey(req.Header.Get("Authorization"))
		if err != nil {
			ErrorOut(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
			return
		}

		ctx := req.Context()
		ctx = context.WithValue(ctx, APIKeyContext("apikey"), key)

		next.ServeHTTP(w, req.WithContext(ctx))
	})
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
