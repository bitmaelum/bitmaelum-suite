package middleware

import (
    "github.com/sirupsen/logrus"
    "net/http"
)

type Tracer struct{}

func (*Tracer) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    logrus.Tracef("%s %s", r.Method, r.URL)

    next.ServeHTTP(w, r)
}
