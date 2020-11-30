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

package dispatcher

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/mborders/artifex"
	"github.com/sirupsen/logrus"
)

var dispatcher *artifex.Dispatcher

// InitDispatcher initializes the dispatcher system with the given number of workers
func InitDispatcher(c int) {
	dispatcher = artifex.NewDispatcher(c, c*10)
	dispatcher.Start()
}

// Dispatch dispatches all active webhooks for the given account and event, with the given payload
func Dispatch(h hash.Hash, evt webhook.EventEnum, payload []byte) error {
	// Nothing to dispatch, as we haven't enabled webhook dispatching
	if dispatcher == nil {
		return nil
	}

	repo := container.Instance.GetWebhookRepo()
	webhooks, err := repo.FetchByHash(h)
	if err != nil {
		return err
	}

	for _, wh := range webhooks {
		if wh.Event != evt {
			continue
		}
		if !wh.Enabled {
			continue
		}

		_ = dispatcher.Dispatch(func() {
			logrus.Debugf("dispatching webhook %s", wh.ID)
			work(wh, payload)
		})
	}

	return nil
}
