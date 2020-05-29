package middleware

import (
    "crypto/subtle"
    "net/http"
)

type BasicAuth struct{}

var username = "user"
var password = "pass"
var realm = "mailv2"

// Basic authentication
func (*BasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    user, pass, ok := r.BasicAuth()

    if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
        w.Header().Set("WWW-Authenticate", `Basic realm="` + realm + `"`)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    next(w, r)
}
