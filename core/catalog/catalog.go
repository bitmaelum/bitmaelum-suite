package catalog

import (
    "github.com/gabriel-vasile/mimetype"
    "github.com/google/uuid"
    "github.com/jaytaph/mailv2/core/account"
    "github.com/jaytaph/mailv2/core/message"
    "io"
    "os"
    "time"
)

type BlockType struct {
    Id          string         `json:"id"`
    Type        string         `json:"type"`
    Size        uint64         `json:"size"`
    Encoding    string         `json:"encoding"`
    Compression string         `json:"compression"`
    Checksum    []message.ChecksumType `json:"checksum"`
    Content     string         `json:"content"`
}

type AttachmentType struct {
    Id          string         `json:"id"`
    MimeType    string         `json:"mimetype"`
    FileName    string         `json:"filename"`
    Size        uint64         `json:"size"`
    Compression string         `json:"compression"`
    Checksum    []message.ChecksumType `json:"checksum"`
}

type Catalog struct {
    From struct {
        Address      string             `json:"address"`
        Name         string             `json:"name"`
        Organisation string             `json:"organisation"`
        ProofOfWork  message.ProofOfWorkType    `json:"proof_of_work"`
        PublicKey    string             `json:"public_key"`
    } `json:"from"`
    To struct {
        Address string `json:"address"`
        Name    string `json:"name"`
    } `json:"to"`
    CreatedAt        time.Time  `json:"created_at"`
    ThreadId         string     `json:"thread_id"`
    Subject          string     `json:"subject"`
    Flags            []string   `json:"flags"`
    Labels           []string   `json:"labels"`
    Catalog struct {
        Blocks []BlockType `json:"blocks"`
        Attachments []AttachmentType `json:"attachments"`
    }
}

type Attachment struct {
    Path        string
    Reader      io.Reader
}

type Block struct {
    Type        string
    Inline      bool
    Content     []byte
}

func NewCatalog(ai *account.AccountInfo) *Catalog {
    c := &Catalog{}

    c.CreatedAt = time.Now()

    c.From.Address = ai.Address
    c.From.Name = ai.Name
    c.From.Organisation = ai.Organisation
    c.From.ProofOfWork.Bits = ai.Pow.Bits
    c.From.ProofOfWork.Proof = ai.Pow.Proof
    c.From.PublicKey = ai.PubKey

    return c;
}

func (c *Catalog) AddBlock(entry Block) error {
    id, err := uuid.NewRandom()
    if err != nil {
        return err
    }

    content := entry.Content
    size := uint64(len(entry.Content))
    if entry.Inline == false {
        stats, err := os.Stat(string(entry.Content))
        if err != nil {
            return err
        }

        content = nil
        size = uint64(stats.Size())
    }

    bt := &BlockType{
        Id:          id.String(),
        Type:        entry.Type,
        Size:        size,
        Encoding:    "",
        Compression: "",
        Checksum:    nil,
        Content:     string(content),
    }

    c.Catalog.Blocks = append(c.Catalog.Blocks, *bt)
    return nil
}

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


    at := &AttachmentType{
        Id:          id.String(),
        MimeType:    mime.String(),
        FileName:    entry.Path,
        Size:        uint64(stats.Size()),
        Compression: "",
        Checksum:    nil,
    }

    c.Catalog.Attachments = append(c.Catalog.Attachments, *at)
    return nil
}