package middleware

import (
    "github.com/sirupsen/logrus"
    "net/http"
)

type Tracer struct{}

// Prints request in log
func (*Tracer) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        logrus.Debugf("%s %s", req.Method, req.URL)

        next.ServeHTTP(w, req)
    })
}
