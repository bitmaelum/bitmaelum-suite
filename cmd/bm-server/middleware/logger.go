package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is a middleware that logs the timing of the given call
type Logger struct{}

// Middleware Logs the request time
func (*Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, req)
		logrus.Tracef("execution time: %s", time.Since(t).String())
	})
}
