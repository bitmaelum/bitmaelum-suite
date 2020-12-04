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
	"encoding/json"
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/google/uuid"
)

var (
	errWebhookNotFound = errors.New("webhook not found")
)

// TypeEnum is the type of the webhook destination
type TypeEnum int

// String convert type to string
func (e TypeEnum) String() string {
	return [...]string{"HTTP"}[e]
}

// Webhook destination types.
const (
	// TypeHTTP
	TypeHTTP TypeEnum = iota // Simple HTTP endpoint
	// TypeAmqp                            // Advanced Message Queue protocol
	// TypeSQS                             // Amazon SQS support
	// TypeSlack                           // Slack support
)

// Type is the webhook structure
type Type struct {
	ID      string    // Id of the webhook (uuidv4)
	Account hash.Hash // account to which this webhook belongs to
	Type    TypeEnum  // type of webhook
	Event   EventEnum // event when this webhook is triggered
	Enabled bool      // true when the webhook is enabled
	Config  string    // config for the given target, json encoded
}

// ConfigHTTP Configuration for TypeHTTP
type ConfigHTTP struct {
	URL string
}

// ConfigSQS Configuration for TypeSQS (not yet implemented)
type ConfigSQS struct {
	Arn             string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

// ConfigSlack Configuration for TypeSlack (not yet implemented)
type ConfigSlack struct {
	WebhookURL string
}

// NewWebhook creates a new webhook
func NewWebhook(account hash.Hash, e EventEnum, t TypeEnum, cfg interface{}) (*Type, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	return &Type{
		ID:      id.String(),
		Account: account,
		Type:    t,
		Event:   e,
		Enabled: false,
		Config:  string(cfgBytes),
	}, nil
}
