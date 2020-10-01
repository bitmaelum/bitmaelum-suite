package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Tracer is a middleware that logs the URL / status code of the call
type Tracer struct {
	http.ResponseWriter
	status int
}

// WriteHeader writes the given header to the response writer
func (t *Tracer) WriteHeader(code int) {
	t.status = code
	t.ResponseWriter.WriteHeader(code)
}

// Middleware Prints request in log
func (*Tracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := Tracer{w, 200}
		next.ServeHTTP(&t, req)
		logrus.Debugf("%s %s (Returned: %d)", req.Method, req.URL, t.status)
	})
}
