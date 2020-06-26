package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type Tracer struct {
	http.ResponseWriter
	status int
}

func (t *Tracer) WriteHeader(code int) {
	t.status = code
	t.ResponseWriter.WriteHeader(code)
}

// Prints request in log
func (*Tracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t := Tracer{w, 200}
		next.ServeHTTP(&t, req)
		logrus.Debugf("%s %s (Returned: %d)", req.Method, req.URL, t.status)
	})
}
