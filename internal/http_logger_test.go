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
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHTTPLogger_LogRequest(t *testing.T) {
	logger, buf := setupLogger()

	body := "hello world"
	req, err := http.NewRequest("get", "https://example.org", strings.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("X-FOO", "bar")
	logger.LogRequest(req)

	assert.Equal(t, buf.String(), "level=trace msg=\"Request get https://example.org\"\nlevel=trace msg=\"HEADER: X-Foo : [bar]\\n\"\nlevel=trace msg=\"log: (11) hello world\"\n")

	buf.Reset()
	req, err = http.NewRequest("get", "https://example.org", http.NoBody)
	assert.NoError(t, err)

	logger.LogRequest(req)
	assert.Equal(t, buf.String(), "level=trace msg=\"Request get https://example.org\"\nlevel=trace msg=\"log: (0) <empty body>\"\n")
}

func TestHTTPLogger_LogResponse(t *testing.T) {
	logger, buf := setupLogger()

	duration := time.Duration(124543242)

	body := "hello world"
	req, err := http.NewRequest("get", "https://example.org", strings.NewReader(body))
	assert.NoError(t, err)

	respBody := " this is the response"
	res := http.Response{
		Status:     "ok",
		StatusCode: 200,
		Header: http.Header{
			"X-Res": []string{"something"},
		},
		Body:          ioutil.NopCloser(strings.NewReader(respBody)),
		ContentLength: int64(len(respBody)),
	}

	logger.LogResponse(req, &res, nil, duration)
	assert.Equal(t, buf.String(), "level=trace msg=\"Response method=get status=200 durationMs=124 https://example.org\"\nlevel=trace msg=\"HEADER: X-Res : [something]\\n\"\n")

	buf.Reset()
	logger.LogResponse(req, &res, errors.New("an error occurred"), duration)
	assert.Equal(t, buf.String(), "level=trace msg=\"log: an error occurred\"\n")
}

func setupLogger() (*HTTPLogger, *bytes.Buffer) {
	logger := NewHTTPLogger()

	formatter := new(logrus.TextFormatter)
	formatter.DisableColors = true
	formatter.DisableTimestamp = true
	logrus.SetFormatter(formatter)

	var b bytes.Buffer
	logrus.SetOutput(&b)

	return logger, &b
}
