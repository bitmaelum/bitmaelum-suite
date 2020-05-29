package catalog

import (
    "github.com/jaytaph/mailv2/core/message"
    "time"
)

//type ChecksumType struct {
//    Hash    string
//    Value   string
//}
//
//type Block struct {
//    Id      string
//    Type    string
//    Size    int64
//    Encoding    string
//    Compression string
//    Checksum    []ChecksumType
//}
//
//type Attachment struct {
//    Id          string
//    Mimetype    string
//    Filename    string
//    Size        int64
//    Encoding    string
//    Compression string
//    Checksum    []ChecksumType
//}
//
//type HeaderType struct {
//    From struct {
//        Email        string
//        Name         string
//        Organisation string
//    }
//    To struct {
//        Email string
//        Name  string
//    }
//    CreatedAt time.Time
//    ThreadId  string
//    Subject   string
//    Labels    []string
//}
//
//type Catalog struct {
//    Header          HeaderType
//    Blocks          []Block
//    Attachments     []Attachment
//}
//
//func CreateCatalog(h *HeaderType, b *[]Block, a *[]Attachment) *Catalog {
//    return &Catalog{
//        Header: *h,
//        Blocks: *b,
//        Attachments: *a,
//    }
//}
//
//func (c *Catalog) Encrypt() *EncryptedCatalog {
//
//}

func CreateCatalog() *message.Catalog {
    cat := &message.Catalog{}

    cat.CreatedAt = time.Now()
    return cat;
}


//func (c *message.Catalog) Encrypt(iv, key string) ([]byte, error) {
//    return []byte{}, nil
//}
