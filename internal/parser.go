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
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// ParseOptions will parse the commandline options given by opts. It will exit when issues arise or help is wanted
func ParseOptions(opts interface{}) {
	parser := flags.NewParser(opts, flags.IgnoreUnknown)
	_, err := parser.Parse()

	if err != nil {
		flagsError, _ := err.(*flags.Error)
		if flagsError.Type == flags.ErrHelp {
			os.Exit(1)
		}

		fmt.Println()
		parser.WriteHelp(os.Stdout)
		fmt.Println()
		os.Exit(1)
	}
}
