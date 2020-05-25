package middleware

import (
    "github.com/sirupsen/logrus"
    "net/http"
    "time"
)

type Logger struct{}

func (*Logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    t := time.Now()

    next.ServeHTTP(w, r)
    logrus.Tracef("execution time: %s \n", time.Now().Sub(t).String())
}
