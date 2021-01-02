// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
