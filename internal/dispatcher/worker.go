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
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/webhook"
)

// work is the main function that will get dispatched as a job. It will do the actual work
func work(w webhook.Type, payload []byte) {
	switch w.Type {
	case webhook.TypeHTTP:
		_ = execHTTP(w, payload)
	}
}

func execHTTP(w webhook.Type, payload []byte) error {
	cfg := &webhook.ConfigHTTP{}
	err := json.Unmarshal([]byte(w.Config), cfg)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post(cfg.URL, "application/json", bytes.NewReader([]byte(payload)))
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("webhook: invalid status code returned from HTTP endpoint")
	}

	return nil
}
