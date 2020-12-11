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
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/spf13/afero"
)

var (
	errBlockFormat              = errors.New("please specify blocks in the format '<type>,<content>' or '<type>,file:<filename>'")
	errBlockNotFound            = errors.New("block not found")
	errAttachmentNotExists      = func(att string) error { return fmt.Errorf("attachment does not exist: %s", att) }
	errAttachmentCannotBeOpened = func(att string) error { return fmt.Errorf("attachment cannot be opened: %s", att) }
)

var fs = afero.NewOsFs()

// BlockType represents a message block as used inside a catalog
type BlockType struct {
	ID          string       `json:"id"`          // BLock identifier UUID
	Type        string       `json:"type"`        // Type of the block. Can be anything message readers can parse.
	Size        uint64       `json:"size"`        // Size of the block in bytes
	Encoding    string       `json:"encoding"`    // Encoding of the block in case it's encoded
	Compression string       `json:"compression"` // Compression used
	Checksum    ChecksumList `json:"checksum"`    // Checksums of the block
	Reader      io.Reader    `json:"-"`           // Reader of the block data
	Key         []byte       `json:"key"`         // Key for decryption
	IV          []byte       `json:"iv"`          // IV for decryption
}

// AttachmentType represents a message attachment as used inside a catalog
type AttachmentType struct {
	ID          string       `json:"id"`          // Attachment identifier UUID
	MimeType    string       `json:"mimetype"`    // Mimetype
	FileName    string       `json:"filename"`    // Filename
	Size        uint64       `json:"size"`        // Size of the attachment in bytes
	Compression string       `json:"compression"` // Compression used
	Checksum    ChecksumList `json:"checksum"`    // Checksums of the data
	Reader      io.Reader    `json:"-"`           // Reader to the attachment data
	Key         []byte       `json:"key"`         // Key for decryption
	IV          []byte       `json:"iv"`          // IV for decryption
}

// Catalog is the structure that represents a message catalog. This will hold all information about the
// actual message, blocks and attachments.
type Catalog struct {
	From struct {
		Address      string                   `json:"address"`       // BitMaelum address of the sender
		Name         string                   `json:"name"`          // Name of the sender
		Organisation string                   `json:"organisation"`  // Organisation of the sender
		ProofOfWork  *proofofwork.ProofOfWork `json:"proof_of_work"` // Sender's proof of work
		PublicKey    *bmcrypto.PubKey         `json:"public_key"`    // Public key of the sender
	} `json:"from"`
	To struct {
		Address string `json:"address"` // Address of the recipient
		Name    string `json:"name"`    // Name of the recipient
	} `json:"to"`
	CreatedAt time.Time `json:"created_at"` // Timestamp when the message was created
	ThreadID  string    `json:"thread_id"`  // Thread ID (and parent ID) in case this message was send in a thread
	Subject   string    `json:"subject"`    // Subject of the message
	Flags     []string  `json:"flags"`      // Flags of the message
	Labels    []string  `json:"labels"`     // Labels for this message

	Blocks      []BlockType      `json:"blocks"`      // Message block info
	Attachments []AttachmentType `json:"attachments"` // Message attachment info
}

// Attachment represents an attachment and reader
type Attachment struct {
	Path   string        // LOCAL path of the attachment. Needed for things like os.Stat()
	Reader io.ReadSeeker // Reader to the attachment file, also needs to seek, as we need to reset after a mimetype check
}

// Block represents a block and reader
type Block struct {
	Type   string    // Type of the block (text, html, default, mobile etc)
	Size   uint64    // Size of the block
	Reader io.Reader // Reader to the block data
}

// NewCatalog initialises a new catalog. This catalog has to be filled with more info, blocks and attachments
func NewCatalog(sender, recipient *address.Address, subject string) *Catalog {
	c := &Catalog{}

	c.CreatedAt = time.Now()
	c.From.Address = sender.String()
	// c.From.Name = info.Name
	// c.From.ProofOfWork = info.Pow
	// c.From.PublicKey = &info.PubKey

	c.Subject = subject
	c.To.Address = recipient.String()

	return c
}

// NewServerCatalog initialises a new (server) catalog. This catalog has to be filled with more info, blocks and attachments
func NewServerCatalog(sender, recipient *hash.Hash, subject string) *Catalog {
	c := &Catalog{}

	c.CreatedAt = time.Now()
	c.From.Address = sender.String()
	c.From.Name = "Postmaster"
	c.From.Address = sender.String()
	c.To.Address = recipient.String()
	c.AddFlags("postmaster")

	c.Subject = subject
	c.To.Address = recipient.String()

	return c
}

// AddFlags adds extra flags to the message
func (c *Catalog) AddFlags(flags ...string) {
	c.Flags = append(c.Flags, flags...)
}

// AddLabels adds extra labels to the message
func (c *Catalog) AddLabels(labels ...string) {
	c.Labels = append(c.Labels, labels...)
}

// SetToAddress sets the address of the recipient
func (c *Catalog) SetToAddress(addr address.Address, fullName string) {
	c.To.Address = addr.String()
	c.To.Name = fullName
}

// AddBlock adds a block to a catalog
func (c *Catalog) AddBlock(entry Block) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var reader = entry.Reader
	var compression = ""

	// Very arbitrary size on when we should compress output first
	if entry.Size >= 1024 {
		reader = ZlibCompress(entry.Reader)
		compression = "zlib"
	}

	// Generate key iv for this block
	iv, key, err := bmcrypto.GenerateIvAndKey()
	if err != nil {
		return err
	}

	// Wrap reader with encryption reader
	reader, err = bmcrypto.GetAesEncryptorReader(iv, key, reader)
	if err != nil {
		return err
	}

	bt := &BlockType{
		ID:          id.String(),
		Type:        entry.Type,
		Size:        entry.Size,
		Encoding:    "",
		Compression: compression,
		Checksum:    nil,
		Reader:      reader,
		Key:         key,
		IV:          iv,
	}

	c.Blocks = append(c.Blocks, *bt)
	return nil
}

// AddAttachment adds an attachment to a catalog
func (c *Catalog) AddAttachment(entry Attachment) error {
	stats, err := fs.Stat(entry.Path)
	if err != nil {
		return err
	}

	mime, err := mimetype.DetectReader(entry.Reader)
	if err != nil {
		return err
	}

	_, _ = entry.Reader.Seek(0, io.SeekStart)

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var reader io.Reader = entry.Reader
	var compression = ""

	// Very arbitrary size on when we should compress output first
	if stats.Size() >= 1024 {
		reader = ZlibCompress(entry.Reader)
		compression = "zlib"
	}

	// Generate Key and IV that we will use for encryption
	iv, key, err := bmcrypto.GenerateIvAndKey()
	if err != nil {
		return err
	}

	// Wrap our reader with the encryption reader
	reader, err = bmcrypto.GetAesEncryptorReader(iv, key, reader)
	if err != nil {
		return err
	}

	at := &AttachmentType{
		ID:          id.String(),
		MimeType:    mime.String(),
		FileName:    path.Base(entry.Path),
		Size:        uint64(stats.Size()),
		Compression: compression,
		Reader:      reader,
		Checksum:    nil, // To be filled in later
		Key:         key,
		IV:          iv,
	}

	c.Attachments = append(c.Attachments, *at)
	return nil
}

// HasBlock returns true when the catalog has the given block type presents
func (c *Catalog) HasBlock(blockType string) bool {
	for _, b := range c.Blocks {
		if b.Type == blockType {
			return true
		}
	}

	return false
}

// GetBlock returns the specified block from the catalog
func (c *Catalog) GetBlock(blockType string) (*BlockType, error) {
	for _, b := range c.Blocks {
		if b.Type == blockType {
			return &b, nil
		}
	}

	return nil, errBlockNotFound
}

// GetFirstBlock returns the first block found in the message
func (c *Catalog) GetFirstBlock() *BlockType {
	return &c.Blocks[0]
}

// GenerateBlocks generates blocks that can be added to a catalog
func GenerateBlocks(b []string) ([]Block, error) {
	// Parse blocks
	var blocks []Block
	for _, block := range b {
		split := strings.SplitN(block, ",", 2)
		if len(split) <= 1 {
			return nil, errBlockFormat
		}

		// By default assume content is inline
		size := int64(len(split[1]))
		var r io.Reader = strings.NewReader(split[1])

		if strings.HasPrefix(split[1], "file:") {
			// Open file as a reader
			f, err := os.Open(strings.TrimPrefix(split[1], "file:"))
			if err != nil {
				return nil, err
			}

			// Read file size
			fi, err := f.Stat()
			if err != nil {
				return nil, err
			}

			r = f
			size = fi.Size()
		}

		blocks = append(blocks, Block{
			Type:   split[0],
			Size:   uint64(size),
			Reader: r,
		})
	}

	return blocks, nil
}

// GenerateAttachments creates message attachments that we can add to a catalog
func GenerateAttachments(a []string) ([]Attachment, error) {
	// Parse attachments
	var attachments []Attachment
	for _, attachment := range a {
		_, err := os.Stat(attachment)
		if os.IsNotExist(err) {
			return nil, errAttachmentNotExists(attachment)
		}

		reader, err := os.Open(attachment)
		if err != nil {
			return nil, errAttachmentCannotBeOpened(attachment)
		}

		attachments = append(attachments, Attachment{
			Path:   attachment,
			Reader: reader,
		})
	}

	return attachments, nil
}
