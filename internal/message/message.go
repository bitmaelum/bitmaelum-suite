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

package message

import (
	"errors"
	"io"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// EncryptedMessage is an encrypted message.
type EncryptedMessage struct {
	ID      string  // Optional message ID (not part of the message)
	Header  *Header // Message header
	Catalog []byte  // Encrypted catalog

	GenerateBlockReader      func(messageID, blockID string) io.Reader      // Generator for block readers
	GenerateAttachmentReader func(messageID, attachmentID string) io.Reader // generator for attachment readers
}

// DecryptedMessage is a message that is fully decrypted and can be read
type DecryptedMessage struct {
	ID      string   // Optional message ID (not part of the message)
	Header  *Header  // Message header (same as in the encrypted message)
	Catalog *Catalog // Decrypted catalog
}

// Decrypt will decrypt the current encrypted message with the given public key and return a decrypted copy
func (em *EncryptedMessage) Decrypt(privKey bmcrypto.PrivKey) (*DecryptedMessage, error) {
	dm := DecryptedMessage{
		ID:     em.ID,
		Header: em.Header,
	}

	// Check signature
	if !VerifyClientHeader(*em.Header) {
		return nil, errors.New("invalid client signature")
	}

	// Decrypt the encryption key
	settings := &bmcrypto.EncryptionSettings{
		Type:          bmcrypto.CryptoType(em.Header.Catalog.Crypto),
		TransactionID: em.Header.Catalog.TransactionID,
	}
	key, err := bmcrypto.Decrypt(privKey, settings, em.Header.Catalog.EncryptedKey)
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
	for idx, blk := range dm.Catalog.Blocks {
		if em.GenerateBlockReader != nil {
			r, err := createReader(blk.IV, blk.Key, blk.Compression, em.GenerateBlockReader(em.ID, blk.ID))
			if err != nil {
				continue
			}
			dm.Catalog.Blocks[idx].Reader = r
		}
	}

	// Add our attachment readers
	for idx, att := range dm.Catalog.Attachments {
		if em.GenerateAttachmentReader != nil {
			r, err := createReader(att.IV, att.Key, att.Compression, em.GenerateAttachmentReader(em.ID, att.ID))
			if err != nil {
				continue
			}
			dm.Catalog.Attachments[idx].Reader = r
		}
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
