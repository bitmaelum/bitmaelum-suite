package message

import (
    "io"
    "time"
)

type BlockType struct {
    Id          string         `json:"id"`
    Type        string         `json:"type"`
    Size        uint64         `json:"size"`
    Encoding    string         `json:"encoding"`
    Compression string         `json:"compression"`
    Checksum    []ChecksumType `json:"checksum"`
    Content     string         `json:"content"`
}

type AttachmentType struct {
    Id          string         `json:"id"`
    MimeType    string         `json:"mimetype"`
    FileName    string         `json:"filename"`
    Size        uint64         `json:"size"`
    Compression string         `json:"compression"`
    Checksum    []ChecksumType `json:"checksum"`
}

type Catalog struct {
    From struct {
        Address      string             `json:"address"`
        Name         string             `json:"name"`
        Organisation string             `json:"organisation"`
        ProofOfWork  ProofOfWorkType    `json:"proof_of_work"`
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
