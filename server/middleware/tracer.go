package middleware

import (
    "github.com/sirupsen/logrus"
    "net/http"
)

type Tracer struct{}

// Prints request in log
func (*Tracer) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logrus.Debugf("%s %s", r.Method, r.URL)

    next.ServeHTTP(w, r)
}
