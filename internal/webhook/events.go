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

package webhook

import (
	"errors"
	"strings"
)

// EventEnum is the event on which the webhook responds
type EventEnum int

// Webhook event types. Everything that a webhook can trigger on
const (
	EventLocalDelivery  EventEnum = iota + 1 // A local delivery has been made (incoming)
	EventFailedDelivery                      // A message is received but not considered correct
	EventRemoteDelivery                      // A message has been send to a remote server (outgoing)
	EventAPIKeyCreated                       // Api Key created
	EventAPIKeyDeleted                       // Api key deleted
	EventAPIKeyUpdated                       // Api key updated
	EventAuthKeyCreated                      // Auth key created
	EventAuthKeyDeleted                      // Auth key deleted
	EventAuthKeyUpdated                      // Auth key updated
	EventWebhookCreated                      // webhook created
	EventWebhookDeleted                      // webhook deleted
	EventWebhookUpdated                      // webhook updated
	EventInviteCreated                       // invite created
	EventAccountCreated                      // account created
	EventTest           EventEnum = 998      // Test event
	EventAll            EventEnum = 999      // All events
)

// EventLabels is the mapping of an event to a string or string to event
var EventLabels = map[string]EventEnum{
	"all":            EventAll,
	"localdelivery":  EventLocalDelivery,
	"faileddelivery": EventFailedDelivery,
	"remotedelivery": EventRemoteDelivery,
	"apikeycreated":  EventAPIKeyCreated,
	"apikeydeleted":  EventAPIKeyDeleted,
	"apikeyupdated":  EventAPIKeyUpdated,
	"authkeycreated": EventAuthKeyCreated,
	"authkeydeleted": EventAuthKeyDeleted,
	"authkeyupdated": EventAuthKeyUpdated,
	"webhookcreated": EventWebhookCreated,
	"webhookdeleted": EventWebhookDeleted,
	"webhookupdated": EventWebhookUpdated,
	"invitecreated":  EventInviteCreated,
	"accountcreated": EventAccountCreated,
	"test":           EventTest,
}

// String convert event to string
func (e EventEnum) String() string {
	for l, ev := range EventLabels {
		if ev == e {
			return l
		}
	}

	return ""
}

// NewEventFromString will create a new event based on the string
func NewEventFromString(s string) (EventEnum, error) {
	for l, ev := range EventLabels {
		if strings.ToLower(s) == l {
			return ev, nil
		}
	}

	return 0, errors.New("unknown event type")
}
