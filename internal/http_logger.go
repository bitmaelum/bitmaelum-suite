package internal

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpLogger struct {
}

func NewHttpLogger() *HttpLogger {
	logrus.SetLevel(logrus.TraceLevel)
	return &HttpLogger{}
}

func (l *HttpLogger) LogRequest(req *http.Request) {
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
		_, err = io.Copy(dest, req.Body)
	}

	req.Body = save

	logrus.Tracef("log: %s", b.String())
}

func (l *HttpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
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
