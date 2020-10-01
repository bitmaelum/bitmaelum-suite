package internal

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPLogger is a structure that allows to log HTTP requests and responses
type HTTPLogger struct {
}

// NewHTTPLogger returns a new http logger
func NewHTTPLogger() *HTTPLogger {
	logrus.SetLevel(logrus.TraceLevel)
	return &HTTPLogger{}
}

// LogRequest will log the request
func (l *HTTPLogger) LogRequest(req *http.Request) {
	var err error
	save := req.Body
	if req.Body != nil {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return
		}
	}

	logrus.Tracef(
		"Request %s %s",
		req.Method,
		req.URL.String(),
	)

	for k, v := range req.Header {
		logrus.Tracef("HEADER: %s : %s\n", k, v)
	}

	var b bytes.Buffer
	if req.Body != nil {
		var dest io.Writer = &b
		_, _ = io.Copy(dest, req.Body)
	}

	req.Body = save

	logrus.Tracef("log: %s", b.String())
}

// LogResponse will log the response
func (l *HTTPLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		logrus.Trace("log: ", err)
	} else {
		logrus.Tracef(
			"Response method=%s status=%d durationMs=%d %s",
			req.Method,
			res.StatusCode,
			duration,
			req.URL.String(),
		)
		for k, v := range res.Header {
			logrus.Tracef("HEADER: %s : %s\n", k, v)
		}

	}
}

// drainBody will duplicate a reader into two separate readers
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
