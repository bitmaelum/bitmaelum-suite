package message

import (
    "io"
    "time"
)

type Header struct {
    From struct {
        Id          string  `json:"id"`
        PublicKey   string  `json:"public_key"`
        ProofOfWork struct {
            Bits    int     `json:"bits"`
            Proof   uint64  `json:"proof"`
        } `json:"proof_of_work"`
    } `json:"from"`
        Id    string    `json:"id"`
    To struct {
    } `json:"to"`
    Catalog struct {
        Size        uint64          `json:"size"`
        Checksum    []ChecksumType  `json:"checksum"`
        Crypto      string          `json:"crypto"`
        Key         []byte          `json:"key"`
        Iv          []byte          `json:"iv"`
    } `json:"catalog"`
}

type Catalog struct {
    From struct {
        Address      string          `json:"address"`
        Name         string          `json:"name"`
        Organisation string          `json:"organisation"`
        ProofOfWork  ProofOfWorkType `json:"proof_of_work"`
        PublicKey    string          `json:"public_key"`
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
        Blocks []struct {
            Id          string         `json:"id"`
            Type        string         `json:"type"`
            Size        uint64         `json:"size"`
            Encoding    string         `json:"encoding"`
            Compression string         `json:"compression"`
            Checksum    []ChecksumType `json:"checksum"`
            Content     string         `json:"content"`
        } `json:"blocks"`
        Attachments []struct {
            Id          string         `json:"id"`
            MimeType    string         `json:"mimetype"`
            FileName    string         `json:"filename"`
            Size        uint64         `json:"size"`
            Compression string         `json:"compression"`
            Checksum    []ChecksumType `json:"checksum"`
        } `json:"attachments"`
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

type ProofOfWorkType struct {
   Bits     int     `json:"bits"`
   Proof    uint64  `json:"proof"`
}

type ChecksumType struct {
   Hash    string  `json:"hash"`
   Value   string  `json:"value"`
}
