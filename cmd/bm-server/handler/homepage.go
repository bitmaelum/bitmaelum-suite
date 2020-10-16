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
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
)

// HomePage Information header on root /
func HomePage(w http.ResponseWriter, req *http.Request) {
	// Simple enough so things like curl work
	htmlVersion := false
	if strings.Contains(req.Header.Get("Accept"), "text/html") {
		htmlVersion = true
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	logo := internal.GetMonochromeASCIILogo()
	if config.Server.Server.VerboseInfo {
		host := fmt.Sprintf("<<< %s >>>", config.Server.Server.Hostname)
		host = fmt.Sprintf("%*s ", (49+len(host))/2, host)
		logo = internal.GetMonochromeASCIILogo() + "\n\n" + host + "\n\n"

		var version bytes.Buffer
		internal.WriteVersionInfo("BitMaelum-Server", &version)
		logo = logo + "\n\n" + version.String()
	}

	if htmlVersion {
		logo = strings.Replace(logo, "\n", "<br>", -1)
		_, _ = w.Write([]byte("<pre>" + logo + "</pre>"))
	} else {
		_, _ = w.Write([]byte(logo))
	}
}
