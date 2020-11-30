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

package mgmt

import (
	"net/http"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/handler"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/internal/httputils"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-server/processor"
	"github.com/bitmaelum/bitmaelum-suite/internal"
)

// FlushQueues handler will flush all the queues normally on tickers
func FlushQueues(w http.ResponseWriter, req *http.Request) {
	k := handler.GetAPIKey(req)
	if !k.HasPermission(internal.PermFlush, nil) {
		httputils.ErrorOut(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Reload configuration and such
	internal.Reload()

	// Flush queues. Note that this means that multiple queue processing can run multiple times
	go processor.ProcessRetryQueue(true)
	go processor.ProcessStuckIncomingMessages()
	go processor.ProcessStuckProcessingMessages()

	_ = httputils.JSONOut(w, http.StatusOK, httputils.StatusOk("Flushing queues"))
}
