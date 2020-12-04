// Copyright (c) 2020 BitMaelum Authors
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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientErrorResponses(t *testing.T) {
	assert.False(t, isErrorResponse([]byte("")))
	assert.False(t, isErrorResponse([]byte("null")))
	assert.False(t, isErrorResponse([]byte("foobar")))
	assert.False(t, isErrorResponse([]byte("{}")))
	assert.False(t, isErrorResponse([]byte("{{{")))
	assert.False(t, isErrorResponse([]byte("{\"error\": \"something\"}")))
	assert.False(t, isErrorResponse([]byte("{\"status\": \"something\"}")))
	assert.False(t, isErrorResponse([]byte("{\"error\": false, \"status\": \"error message\"}")))
	assert.False(t, isErrorResponse([]byte("{\"error\": false, \"status\": \"\"}")))
	assert.True(t, isErrorResponse([]byte("{\"error\": true, \"status\": \"error message\"}")))
	assert.True(t, isErrorResponse([]byte("{\"error\": true}")))

	assert.NoError(t, GetErrorFromResponse([]byte("")))
	assert.NoError(t, GetErrorFromResponse([]byte("null")))
	assert.NoError(t, GetErrorFromResponse([]byte("foobar")))
	assert.NoError(t, GetErrorFromResponse([]byte("{}")))
	assert.NoError(t, GetErrorFromResponse([]byte("{{{")))
	assert.NoError(t, GetErrorFromResponse([]byte("{\"error\": \"something\"}")))
	assert.NoError(t, GetErrorFromResponse([]byte("{\"status\": \"something\"}")))
	assert.NoError(t, GetErrorFromResponse([]byte("{\"error\": false, \"status\": \"error message\"}")))
	assert.NoError(t, GetErrorFromResponse([]byte("{\"error\": false, \"status\": \"\"}")))
	assert.Errorf(t, GetErrorFromResponse([]byte("{\"error\": true}")), "unnown erorr")
	assert.Errorf(t, GetErrorFromResponse([]byte("{\"error\": true, \"status\": \"error message\"}")), "error message")
}

func TestCanonicalHost(t *testing.T) {
	assert.Equal(t, "https://foo.example.org:2424", CanonicalHost("foo.example.org"))
	assert.Equal(t, "https://foo.example.org:80", CanonicalHost("foo.example.org:80"))
	assert.Equal(t, "http://foo.example.org:2424", CanonicalHost("http://foo.example.org"))
	assert.Equal(t, "mail://foo.example.org:2424", CanonicalHost("mail://foo.example.org"))
	assert.Equal(t, "mail://foo.example.org:1234", CanonicalHost("mail://foo.example.org:1234"))
	assert.Equal(t, "mail://192.168.1.1:1234", CanonicalHost("mail://192.168.1.1:1234"))
	assert.Equal(t, "https://192.168.1.1:1234", CanonicalHost("192.168.1.1:1234"))
	assert.Equal(t, "https://192.168.1.1:2424", CanonicalHost("192.168.1.1"))
	assert.Equal(t, "https://[::1]:80", CanonicalHost("[::1]:80"))
	assert.Equal(t, "https://[::1]:2424", CanonicalHost("[::1]"))
	assert.Equal(t, "http://[::1]:2424", CanonicalHost("http://[::1]"))
	assert.Equal(t, "http://[2a04:c44:e00:147a:441:aaff:fe00:1d2]:8080", CanonicalHost("http://[2a04:c44:e00:147a:441:aaff:fe00:1d2]:8080"))
	assert.Equal(t, "https://[2a04:c44:e00:147a:441:aaff:fe00:1d2]:8080", CanonicalHost("https://[2a04:c44:e00:147a:441:aaff:fe00:1d2]:8080"))
	assert.Equal(t, "https://[2a04:c44:e00:147a:441:aaff:fe00:1d2]:2424", CanonicalHost("[2a04:c44:e00:147a:441:aaff:fe00:1d2]"))
}

// // GetErrorFromResponse will return an error generated from the body
// func GetErrorFromResponse(body []byte) error {
// 	type errorStatus struct {
// 		Error  bool   `json:"error"`
// 		Status string `json:"status"`
// 	}
//
// 	s := &errorStatus{Error: true}
// 	err := json.Unmarshal(body, s)
// 	if err != nil {
// 		return errNoSuccess
// 	}
//
// 	return errors.New(s.Status)
// }
