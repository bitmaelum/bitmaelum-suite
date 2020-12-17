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
	"errors"
	"io"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// EncryptedMessage is an encrypted message.
type EncryptedMessage struct {
	BoxID   string  // Optional block ID where this message resied (not part of the message)
	ID      string  // Optional message ID (not part of the message)
	Header  *Header // Message header
	Catalog []byte  // Encrypted catalog

	GenerateBlockReader      func(boxID, messageID, blockID string) io.Reader      // Generator for block readers
	GenerateAttachmentReader func(boxID, messageID, attachmentID string) io.Reader // generator for attachment readers
}

// DecryptedMessage is a message that is fully decrypted and can be read
type DecryptedMessage struct {
	ID                string               // Optional message ID (not part of the message)
	Header            *Header              // Message header (same as in the encrypted message)
	Catalog           *Catalog             // Decrypted catalog
	BlockReaders      map[string]io.Reader // Readers to the (decrypted and decompressed) blocks
	AttachmentReaders map[string]io.Reader // Readers to the (decrypted and decompressed) attachments
}

// Decrypt will decrypt the current encrypted message with the given public key and return a decrypted copy
func (em *EncryptedMessage) Decrypt(privKey bmcrypto.PrivKey) (*DecryptedMessage, error) {
	dm := DecryptedMessage{
		Header:            em.Header,
		BlockReaders:      make(map[string]io.Reader),
		AttachmentReaders: make(map[string]io.Reader),
	}

	// Check signature
	if !VerifyClientHeader(*em.Header) {
		return nil, errors.New("invalid client signature")
	}

	// Decrypt the encryption key
	key, err := bmcrypto.Decrypt(privKey, em.Header.Catalog.TransactionID, em.Header.Catalog.EncryptedKey)
	if err != nil {
		return nil, err
	}

	// Decrypt the catalog
	dm.Catalog = &Catalog{}
	err = bmcrypto.CatalogDecrypt(key, em.Catalog, dm.Catalog)
	if err != nil {
		return nil, errors.New("cannot decrypt")
	}

	// Add our block readers
	for _, blk := range dm.Catalog.Blocks {
		r, err := createReader(blk.IV, blk.Key, blk.Compression, em.GenerateBlockReader(em.BoxID, em.ID, blk.ID))
		if err != nil {
			continue
		}
		dm.BlockReaders[blk.ID] = r
	}

	// Add our attachment readers
	for _, att := range dm.Catalog.Attachments {
		r, err := createReader(att.IV, att.Key, att.Compression, em.GenerateAttachmentReader(em.BoxID, em.ID, att.ID))
		if err != nil {
			continue
		}
		dm.AttachmentReaders[att.ID] = r
	}

	return &dm, nil
}

func createReader(iv []byte, key []byte, compression string, reader io.Reader) (io.Reader, error) {
	r, err := bmcrypto.GetAesDecryptorReader(iv, key, reader)
	if err != nil {
		return nil, err
	}

	switch compression {
	case "zlib":
		r, err = ZlibDecompress(r)
		if err != nil {
			return nil, err
		}
	default:
		// do nothing
	}

	return r, nil
}
