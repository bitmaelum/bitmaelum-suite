// Copyright (c) 2021 BitMaelum Authors
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

package bitmaelumClient

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/pkg/errors"
)

func (b *BitMaelumClient) SendSimpleMessage(to, subject, body string) error {
	senderInfo, _ := b.resolverService.ResolveAddress(b.client.Address.Hash())

	// Check recipient
	toAddr, err := address.NewAddress(to)
	if err != nil {
		return errors.Wrap(err, "check recipient address")
	}

	recipientInfo, err := b.resolverService.ResolveAddress(toAddr.Hash())
	if err != nil {
		return errors.Wrap(err, "resolve recipient address")
	}

	// Setup addressing
	addressing := message.NewAddressing(message.SignedByTypeOrigin)
	addressing.AddSender(b.client.Address, nil, b.client.Name, *b.client.PrivateKey, senderInfo.RoutingInfo.Routing)
	addressing.AddRecipient(toAddr, nil, &recipientInfo.PublicKey)

	// Setup blocks
	var blocks []string
	blocks = append(blocks, "default,"+body)

	// Compose mail
	envelope, err := message.Compose(addressing, subject, blocks, nil)
	if err != nil {
		return errors.Wrap(err, "composing mail")
	}

	// Send mail
	client, err := api.NewAuthenticated(*b.client.Address, *b.client.PrivateKey, senderInfo.RoutingInfo.Routing, nil)
	if err != nil {
		return errors.Wrap(err, "setting api")
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		return errors.Wrap(err, "sending message")
	}

	return nil
}

/*
func (b *BitMaelumClient) SendSimpleMessage(fromAcc, fromName, to, privKey, subject, body string) error {
	config.Client.Resolver.Remote.Enabled = true
	config.Client.Resolver.Remote.URL = b.resolverURL

	svc := container.Instance.GetResolveService()

	// Check sender
	fromAddr, err := address.NewAddress(fromAcc)
	if err != nil {
		return err
	}

	senderInfo, _ := svc.ResolveAddress(fromAddr.Hash())

	// Check recipient
	toAddr, err := address.NewAddress(to)
	if err != nil {
		return err
	}

	recipientInfo, err := svc.ResolveAddress(toAddr.Hash())
	if err != nil {
		return err
	}

	// Convert privKey string to bmcrypto
	pk, err := bmcrypto.PrivateKeyFromString(privKey)
	if err != nil {
		return err
	}

	// Setup addressing
	addressing := message.NewAddressing(message.SignedByTypeOrigin)
	addressing.AddSender(fromAddr, nil, fromName, *pk, senderInfo.RoutingInfo.Routing)
	addressing.AddRecipient(toAddr, nil, &recipientInfo.PublicKey)

	// Setup blocks
	var blocks []string
	blocks = append(blocks, "default,"+body)

	// Compose mail
	envelope, err := message.Compose(addressing, subject, blocks, nil)
	if err != nil {
		return err
	}

	// Send mail
	client, err := api.NewAuthenticated(*fromAddr, *pk, senderInfo.RoutingInfo.Routing, nil)
	if err != nil {
		return err
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		return err
	}

	return nil
}
*/
