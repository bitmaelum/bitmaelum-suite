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

package dispatcher

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
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

// dispatch dispatches all active webhooks for the given account and event, with the given payload
func dispatch(h hash.Hash, evt webhook.EventEnum, payload map[string]interface{}) error {
	logrus.Debugf("dispatching for %s (%s)", h, evt)

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
		// Enabled webhooks only
		if !wh.Enabled {
			continue
		}

		// Make sure we match either the event or webhook acts on "all"
		if wh.Event != evt && wh.Event != webhook.EventAll {
			continue
		}

		// Add meta data to payload
		payload["meta"] = map[string]string{
			"id":      wh.ID,
			"account": wh.Account.String(),
			"event":   evt.String(), // We need to use the event, otherwise we might end up with event "all"
		}

		wh := wh
		_ = dispatcher.Dispatch(func() {
			logrus.Tracef("dispatching webhook %s", wh.ID)
			Work(wh, payload)
		})
	}

	return nil
}

// DispatchRemoteDelivery dispatches event
func DispatchRemoteDelivery(h hash.Hash, header *message.Header, msgID string) error {
	payload := map[string]interface{}{
		"message": map[string]string{
			"from": header.From.Addr.String(),
			"to":   header.To.Addr.String(),
			"id":   msgID,
		},
	}

	return dispatch(h, webhook.EventRemoteDelivery, payload)
}

// DispatchLocalDelivery dispatches event
func DispatchLocalDelivery(h hash.Hash, header *message.Header, msgID string) error {
	payload := map[string]interface{}{
		"message": map[string]string{
			"from": header.From.Addr.String(),
			"to":   header.To.Addr.String(),
			"id":   msgID,
		},
	}

	return dispatch(h, webhook.EventLocalDelivery, payload)
}

// DispatchAPIKeyCreate dispatches event
func DispatchAPIKeyCreate(h hash.Hash, k key.APIKeyType) error {
	payload := map[string]interface{}{
		"key": k,
	}

	return dispatch(h, webhook.EventAPIKeyCreated, payload)
}

// DispatchAPIKeyDelete dispatches event
func DispatchAPIKeyDelete(h hash.Hash, k key.APIKeyType) error {
	payload := map[string]interface{}{
		"key": map[string]string{
			"id": k.ID,
		},
	}

	return dispatch(h, webhook.EventAPIKeyDeleted, payload)
}

// DispatchAuthKeyCreate dispatches event
func DispatchAuthKeyCreate(h hash.Hash, k key.AuthKeyType) error {
	payload := map[string]interface{}{
		"key": k,
	}

	return dispatch(h, webhook.EventAuthKeyCreated, payload)
}

// DispatchAuthKeyDelete dispatches event
func DispatchAuthKeyDelete(h hash.Hash, k key.AuthKeyType) error {
	payload := map[string]interface{}{
		"key": map[string]string{
			"id": k.Fingerprint,
		},
	}

	return dispatch(h, webhook.EventAuthKeyDeleted, payload)
}

// DispatchWebhookCreate dispatches event
func DispatchWebhookCreate(h hash.Hash, wh webhook.Type) error {
	payload := map[string]interface{}{
		"webhook": wh,
	}

	return dispatch(h, webhook.EventWebhookCreated, payload)
}

// DispatchWebhookUpdate dispatches event
func DispatchWebhookUpdate(h hash.Hash, wh webhook.Type) error {
	payload := map[string]interface{}{
		"webhook": wh,
	}

	return dispatch(h, webhook.EventWebhookUpdated, payload)
}

// DispatchWebhookDelete dispatches event
func DispatchWebhookDelete(h hash.Hash, wh webhook.Type) error {
	payload := map[string]interface{}{
		"webhook": map[string]string{
			"id": wh.ID,
		},
	}

	return dispatch(h, webhook.EventAuthKeyDeleted, payload)
}
