package message

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"io"
	"os"
	"time"
)

// BlockType represents a message block as used inside a catalog
type BlockType struct {
	ID          string     `json:"id"`          // BLock identifier UUID
	Type        string     `json:"type"`        // Type of the block. Can be anything message readers can parse.
	Size        uint64     `json:"size"`        // Size of the block in bytes
	Encoding    string     `json:"encoding"`    // Encoding of the block in case it's encoded
	Compression string     `json:"compression"` // Compression used
	Checksum    []Checksum `json:"checksum"`    // Checksums of the block
	Reader      io.Reader  `json:"-"`           // Reader of the block data
	Key         []byte     `json:"key"`         // Key for decryption
	IV          []byte     `json:"iv"`          // IV for decryption
}

// AttachmentType represents a message attachment as used inside a catalog
type AttachmentType struct {
	ID          string     `json:"id"`          // Attachment identifier UUID
	MimeType    string     `json:"mimetype"`    // Mimetype
	FileName    string     `json:"filename"`    // Filename
	Size        uint64     `json:"size"`        // Size of the attachment in bytes
	Compression string     `json:"compression"` // Compression used
	Checksum    []Checksum `json:"checksum"`    // Checksums of the data
	Reader      io.Reader  `json:"-"`           // Reader to the attachment data
	Key         []byte     `json:"key"`         // Key for decryption
	IV          []byte     `json:"iv"`          // IV for decryption
}

// Catalog is the structure that represents a message catalog. This will hold all information about the
// actual message, blocks and attachments.
type Catalog struct {
	From struct {
		Address      string           `json:"address"`       // BitMaelum address of the sender
		Name         string           `json:"name"`          // Name of the sender
		Organisation string           `json:"organisation"`  // Organisation of the sender
		ProofOfWork  core.ProofOfWork `json:"proof_of_work"` // Sender's proof of work
		PublicKey    string           `json:"public_key"`    // Public key of the sender
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
	Path   string    // LOCAL path of the attachment. Needed for things like os.Stat()
	Reader io.Reader // Reader to the attachment file
}

// Block represents a block and reader
type Block struct {
	Type   string    // Type of the block (text, html, default, mobile etc)
	Size   uint64    // Size of the block
	Reader io.Reader // Reader to the block data
}

// NewCatalog initialises a new catalog. This catalog has to be filled with more info, blocks and attachments
func NewCatalog(ai *core.AccountInfo) *Catalog {
	c := &Catalog{}

	c.CreatedAt = time.Now()

	c.From.Address = ai.Address
	c.From.Name = ai.Name
	c.From.Organisation = ai.Organisation
	c.From.ProofOfWork.Bits = ai.Pow.Bits
	c.From.ProofOfWork.Proof = ai.Pow.Proof
	c.From.PublicKey = string(ai.PubKey)

	return c
}

// AddBlock adds a block to a catalog
func (c *Catalog) AddBlock(entry Block) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var reader io.Reader = entry.Reader
	var compression = ""

	// Very arbitrary size on when we should compress output first
	if entry.Size >= 1024 {
		reader = core.ZlibCompress(entry.Reader)
		compression = "zlib"
	}

	// Generate key iv for this block
	iv, key, err := GenerateIvAndKey()
	if err != nil {
		return err
	}

	// Wrap reader with encryption reader
	reader, err = GetAesEncryptorReader(iv, key, reader)
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
	stats, err := os.Stat(entry.Path)
	if err != nil {
		return err
	}

	mime, err := mimetype.DetectReader(entry.Reader)
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	var reader io.Reader = entry.Reader
	var compression = ""

	// Very arbitrary size on when we should compress output first
	if stats.Size() >= 1024 {
		reader = core.ZlibCompress(entry.Reader)
		compression = "zlib"
	}

	// Generate Key and IV that we will use for encryption
	iv, key, err := GenerateIvAndKey()
	if err != nil {
		return err
	}

	// Wrap our reader with the encryption reader
	reader, err = GetAesEncryptorReader(iv, key, reader)
	if err != nil {
		return err
	}

	at := &AttachmentType{
		ID:          id.String(),
		MimeType:    mime.String(),
		FileName:    entry.Path,
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

// GenerateIvAndKey generate a random IV and key
// @TODO: This is not a good spot. We should store it in the encrypt page, but this gives us a import cycle
func GenerateIvAndKey() ([]byte, []byte, error) {
	iv := make([]byte, 16)
	n, err := rand.Read(iv)
	if n != 16 || err != nil {
		return nil, nil, err
	}

	key := make([]byte, 32)
	n, err = rand.Read(key)
	if n != 32 || err != nil {
		return nil, nil, err
	}

	return iv, key, nil
}

// GetAesEncryptorReader returns a reader that automatically encrypts reader blocks through CFB stream
// @TODO: This is not a good spot. We should store it in the encrypt page, but this gives us a import cycle
func GetAesEncryptorReader(iv []byte, key []byte, r io.Reader) (io.Reader, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, err
}

// GetAesDecryptorReader returns a reader that automatically decrypts reader blocks through CFB stream
func GetAesDecryptorReader(iv []byte, key []byte, r io.Reader) (io.Reader, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, err
}
