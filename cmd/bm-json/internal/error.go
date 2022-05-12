// Copyright (c) 2022 BitMaelum Authors
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
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/output"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
)

// JwtJSONErrorFunc is the generic error handler that will catch any timing issues with the clock
func JwtJSONErrorFunc(_ *http.Request, resp *http.Response) {
	if resp.StatusCode != 401 {
		return
	}

	// Read body
	b, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	err := api.GetErrorFromResponse(b)
	if err != nil && err.Error() == "token time not valid" {
		output.JSONErrorStrOut("JWT token time mismatch")
		os.Exit(1)
	}

	// Whoops.. not an error. Let's pretend nothing happened and create a new buffer so we can read the body again
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
}
