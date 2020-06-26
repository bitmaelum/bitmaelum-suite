package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Logger struct{}

// Logs the request time
func (*Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, req)
		logrus.Tracef("execution time: %s \n", time.Now().Sub(t).String())
	})
}
