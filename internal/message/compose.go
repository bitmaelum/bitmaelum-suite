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

package message

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// Compose will create a new message and places it inside an envelope. This can be used for actual sending the message
func Compose(addressing Addressing, subject string, b, a []string) (*Envelope, error) {
	cat, err := generateCatalog(addressing.Sender.Address, addressing.Recipient.Address, subject, b, a)
	if err != nil {
		return nil, err
	}

	header, err := generateHeader(addressing.Sender.Address.Hash(), addressing.Recipient.Address.Hash(), SignedByTypeOrigin)
	if err != nil {
		return nil, err
	}

	envelope, err := NewEnvelope()
	if err != nil {
		return nil, err
	}

	err = envelope.AddCatalog(cat)
	if err != nil {
		return nil, err
	}

	err = envelope.AddHeader(header)
	if err != nil {
		return nil, err
	}

	// Close the envelope for sending
	err = envelope.CloseAndEncrypt(addressing.Sender.PrivKey, addressing.Recipient.PubKey)
	if err != nil {
		return nil, err
	}

	return envelope, nil
}

// ServerCompose will create a new (server) message and places it inside an envelope.
func ServerCompose(sender, recipient hash.Hash, senderPrivKey *bmcrypto.PrivKey, recipientPublicKey *bmcrypto.PubKey, subject string, b, a []string) (*Envelope, error) {
	cat, err := generateServerCatalog(sender, recipient, subject, b, a)
	if err != nil {
		return nil, err
	}

	header, err := generateHeader(sender, recipient, SignedByTypeServer)
	if err != nil {
		return nil, err
	}

	envelope, err := NewEnvelope()
	if err != nil {
		return nil, err
	}

	err = envelope.AddCatalog(cat)
	if err != nil {
		return nil, err
	}

	err = envelope.AddHeader(header)
	if err != nil {
		return nil, err
	}

	// Close the envelope for sending
	err = envelope.CloseAndEncrypt(senderPrivKey, recipientPublicKey)
	if err != nil {
		return nil, err
	}

	return envelope, nil
}

// Generate a header file based on the info provided
func generateHeader(sender, recipient hash.Hash, origin SignedByType) (*Header, error) {
	header := &Header{}

	header.From.Addr = sender
	header.To.Addr = recipient
	header.From.SignedBy = origin

	return header, nil
}

// Generate a catalog filled with blocks and attachments
func generateCatalog(sender, recipient address.Address, subject string, b, a []string) (*Catalog, error) {
	// Create a new catalog
	cat := NewCatalog(&sender, &recipient, subject)

	return finishCatalog(cat, b, a)
}

// Generate a catalog filled with blocks and attachments
func generateServerCatalog(sender, recipient hash.Hash, subject string, b, a []string) (*Catalog, error) {
	// Create a new catalog
	cat := NewServerCatalog(&sender, &recipient, subject)

	return finishCatalog(cat, b, a)
}

func finishCatalog(cat *Catalog, b, a []string) (*Catalog, error) {
	// Add blocks to catalog
	blocks, err := GenerateBlocks(b)
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		err := cat.AddBlock(block)
		if err != nil {
			return nil, err
		}
	}

	// Add attachments to catalog
	attachments, err := GenerateAttachments(a)
	if err != nil {
		return nil, err
	}
	for _, attachment := range attachments {
		err := cat.AddAttachment(attachment)
		if err != nil {
			return nil, err
		}
	}

	return cat, nil
}
