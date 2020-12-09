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

package output

import (
	"encoding/json"
	"errors"
	"fmt"
)

// JSONT is a simple type for generating JSON output
type JSONT map[string]interface{}

// JSONOut outputs a specific interface
func JSONOut(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return
	}

	fmt.Print(string(b))
}

// JSONErrorOut outputs an error
func JSONErrorOut(err error) {
	v := map[string]interface{}{
		"error": err.Error(),
	}

	JSONOut(v)
}

// JSONErrorStrOut outputs an error string
func JSONErrorStrOut(s string) {
	JSONErrorOut(errors.New(s))
}
