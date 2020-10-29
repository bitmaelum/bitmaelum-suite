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

package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/middleware/auth"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/sirupsen/logrus"
)

// OutputResponse is a generic output response
type OutputResponse struct {
	Error  bool   `json:"error,omitempty"`
	Status string `json:"status"`
}

// JSONOut outputs the given data structure to JSON
func JSONOut(w http.ResponseWriter, v interface{}) error {
	var data []byte
	var err error

	// Indent or not, depends on the x-pretty-json header set in the response. This is set by the prettyjson middleware
	if w.Header().Get("x-pretty-json") == "1" {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}
	w.Header().Del("x-pretty-json")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(StatusError("Malformed JSON: " + err.Error()))
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	return err
}

// ErrorOut outputs an error
func ErrorOut(w http.ResponseWriter, code int, msg string) {
	logrus.Debugf("Returning error (%d): %s", code, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(StatusError(msg))
}

// StatusOk Return an ok status response
func StatusOk(status string) *OutputResponse {
	return &OutputResponse{
		Status: status,
	}
}

// StatusError Return an error status response
func StatusError(status string) *OutputResponse {
	return &OutputResponse{
		Error:  true,
		Status: status,
	}
}

// DecodeBody decodes a JSON body or write/return an error if an error occurs
func DecodeBody(w http.ResponseWriter, body io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(StatusError("Malformed JSON: " + err.Error()))
		return err
	}

	return nil
}

// GetAPIKey returns the api key stored in the request context. If not found, it will return a dummy key with no permissions
func GetAPIKey(req *http.Request) *key.APIKeyType {
	val := req.Context().Value(auth.APIKeyContext)
	if val == nil {
		return &key.APIKeyType{}
	}

	return val.(*key.APIKeyType)
}

// IsAPIKeyAuthenticated returns true when the given request is authenticated by a api key
func IsAPIKeyAuthenticated(req *http.Request) bool {
	return req.Context().Value("auth_method") == "*middleware.APIKey"
}
