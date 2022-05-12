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

package handlers

import (
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
)

// ComposeMessage composes a new message from the given account Info to the "to" with given subject, blocks and attachments
func ComposeMessage(addressing message.Addressing, subject string, b, a []string) error {
	envelope, err := message.Compose(addressing, subject, b, a)
	if err != nil {
		return err
	}

	// Setup API connection to the server
	client, err := api.NewAuthenticated(*addressing.Sender.Address, *addressing.Sender.PrivKey, addressing.Sender.Host, internal.JwtErrorFunc)
	if err != nil {
		return err
	}

	// and finally send
	err = messages.Send(*client, envelope)
	if err != nil {
		return err
	}
	return nil
}
