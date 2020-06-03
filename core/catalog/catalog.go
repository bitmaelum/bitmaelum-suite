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
