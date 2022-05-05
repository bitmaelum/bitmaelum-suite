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
// FOR A PARTICULeAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package message

import (
	"bytes"
	"errors"
	"io"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

var errAlreadyClosed = errors.New("envelope is already closed and encrypted")

// Envelope is a simple structure that will keep a header and catalog together.
type Envelope struct {
	Header            *Header               // The message header
	EncryptedCatalog  []byte                // The catalog in encrypted bytes
	BlockReaders      map[string]*io.Reader // Readers for the blocks
	AttachmentReaders map[string]*io.Reader // Readers for the attachments

	closed     bool     // True when the envelope is closed (and encrypted)
	catalog    *Catalog // The opened catalog
	catalogKey []byte   // The key used for encryption
}

// Addressing is the configuration for an envelope (@TODO: bad naming)
type Addressing struct {
	Sender struct {
		Address *address.Address  // BitMaelum Address
		Hash    *hash.Hash        // Optional address hash if the address is unknown (for instance, server messages). If the address is set, this field is ignored
		Name    string            // Additional name (if known)
		PrivKey *bmcrypto.PrivKey // Private key of the sender
		Host    string            // Host address of the sender mail server
	}
	Recipient struct {
		Address *address.Address // BitMaelum Address
		Hash    *hash.Hash       // Optional address hash if the address is unknown (for instance, server messages). If the address is set, this field is ignored
		PubKey  *bmcrypto.PubKey // Public key
	}
	AuthorizedBy struct {
		PubKey    *bmcrypto.PubKey // Public key of the authorized user
		Signature string           // signature of the signed public key by the sender.address
	}
	Type SignedByType // Type of the addressing
}

// NewAddressing sets up a new addressing struct that can be used for composing and sending a message. Use the fluent builder
// methods to populate the sender and the receiver
func NewAddressing(signType SignedByType) Addressing {
	return Addressing{
		Type: signType,
	}
}

// AddSender will add sender information to the addressing
func (a *Addressing) AddSender(addr *address.Address, h *hash.Hash, name string, key bmcrypto.PrivKey, host string) {
	a.Sender.Address = addr
	a.Sender.Hash = h
	a.Sender.Name = name
	a.Sender.PrivKey = &key
	a.Sender.Host = host
}

// AddRecipient will add recipient information to the addressing
func (a *Addressing) AddRecipient(addr *address.Address, h *hash.Hash, key *bmcrypto.PubKey) {
	a.Recipient.Address = addr
	a.Recipient.Hash = h
	a.Recipient.PubKey = key
}

// NewEnvelope creates a new (open) envelope which is used for holding a complete message
func NewEnvelope() (*Envelope, error) {
	var err error

	envelope := &Envelope{
		closed:            false,
		BlockReaders:      make(map[string]*io.Reader),
		AttachmentReaders: make(map[string]*io.Reader),
	}

	// Always create a catalog key for later use
	envelope.catalogKey, err = bmcrypto.CreateCatalogKey()
	if err != nil {
		return nil, err
	}

	return envelope, nil
}

// AddHeader will add a header to the envelope
func (e *Envelope) AddHeader(hdr *Header) error {
	if e.closed {
		return errAlreadyClosed
	}

	e.Header = hdr
	return nil
}

// AddCatalog will add a catalog to the envelope
func (e *Envelope) AddCatalog(cat *Catalog) error {
	if e.closed {
		return errAlreadyClosed
	}

	e.catalog = cat

	// Remove all readers first, then add all the new ones
	e.BlockReaders = make(map[string]*io.Reader)
	e.AttachmentReaders = make(map[string]*io.Reader)

	for i := range cat.Blocks {
		block := cat.Blocks[i]
		e.BlockReaders[block.ID] = &block.Reader
	}

	for i := range cat.Attachments {
		att := cat.Attachments[i]
		e.AttachmentReaders[att.ID] = &att.Reader
	}

	return nil
}

// CloseAndEncrypt will close an envelope, and make sure all settings are set correctly for sending the message
func (e *Envelope) CloseAndEncrypt(senderPrivKey *bmcrypto.PrivKey, recipientPubKey *bmcrypto.PubKey) error {
	var err error

	if e.closed {
		return errAlreadyClosed
	}

	// Encrypt the catalog
	e.EncryptedCatalog, err = bmcrypto.JSONEncrypt(e.catalogKey, e.catalog)
	if err != nil {
		return err
	}

	// Calculate checksums of the encrypted catalog
	r := bytes.NewBuffer(e.EncryptedCatalog)
	e.Header.Catalog.Checksum, err = CalculateChecksums(r)
	if err != nil {
		return err
	}

	// Set catalog information in the header
	e.Header.Catalog.Size = uint64(len(e.EncryptedCatalog))
	ek, settings, err := bmcrypto.Encrypt(*recipientPubKey, e.catalogKey)
	if err != nil {
		return err
	}
	e.Header.Catalog.Crypto = string(settings.Type)
	e.Header.Catalog.EncryptedKey = ek
	e.Header.Catalog.TransactionID = settings.TransactionID

	// Sign the header
	err = SignClientHeader(e.Header, *senderPrivKey)
	if err != nil {
		return err
	}

	// All done. Close the envelope, and remove the open catalog
	e.catalog = nil
	e.closed = true
	return nil
}
