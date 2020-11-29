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

package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/google/uuid"
)

var (
	errWebhookNotFound = errors.New("webhook not found")
)

type TypeEnum int
type EventEnum int

const (
	TypeHTTP TypeEnum = iota // Simple HTTP endpoint
	// TypeAmqp                            // Advanced Message Queue protocol
	// TypeSQS                             // Amazon SQS support
	// TypeSlack                           // Slack support

	EventNewMessage EventEnum = iota
)

type Type struct {
	ID      string    // Id of the webhook
	Account hash.Hash // account to which this webhook belongs to
	Type    TypeEnum  // type of webhook
	Event   EventEnum // event when this webhook is triggered
	Enabled bool      // true when the webook is enabled
	Config  []byte    // config for the given target.
}

type ConfigHTTP struct {
	Url string
}

type ConfigSQS struct {
	Arn             string
	Region          string
	AccesKeyId      string
	SecretAccessKey string
}

type ConfigSlack struct {
	WebhookUrl string
}

// @TODO: should we just have a queue system that will run at max 10 go routines for instance? Otherwise all go threads
// might be depleted.

// Execute will run the current webhook inside a separate go routine
func (w *Type) Execute(payload string) {
	go func() {
		err := w.executeWebhook(payload)

		if err != nil {
			// @TODO: what to do when the webhook fails?
		}
	}()
}

func (w *Type) executeWebhook(payload string) error {
	switch w.Type {
	case TypeHTTP:
		return w.execHttp(payload)
	}

	return errors.New("webhook: type not supported")
}

func (w *Type) execHttp(payload string) error {
	cfg := &ConfigHTTP{}
	err := json.Unmarshal(w.Config, cfg)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Post(cfg.Url, "application/json", bytes.NewReader([]byte(payload)))
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return errors.New("webhook: invalid status code returned from HTTP endpoint")
	}

	return nil
}

func NewWebhook(account hash.Hash, t TypeEnum, e EventEnum, cfg []byte) (*Type, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Type{
		ID:      id.String(),
		Account: account,
		Type:    t,
		Event:   e,
		Enabled: false,
		Config:  cfg,
	}, nil
}
