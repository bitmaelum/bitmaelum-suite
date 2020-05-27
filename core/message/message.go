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
            Proof   int64   `json:"proof"`
        } `json:"proof_of_work"`
    } `json:"from"`
        Id    string    `json:"id"`
    To struct {
    } `json:"to"`
    Catalog struct {
        Size        int64           `json:"size"`
        Checksum    []ChecksumType  `json:"checksum"`
        Crypto      string          `json:"crypto"`
        Key         []byte          `json:"key"`
        Iv          []byte          `json:"iv"`
    } `json:"catalog"`
}

type Catalog struct {
    From struct {
        Email        string          `json:"email"`
        Name         string          `json:"name"`
        Organisation string          `json:"organisation"`
        ProofOfWork  ProofOfWorkType `json:"proof_of_work"`
        PublicKey    string          `json:"public_key"`
    } `json:"from"`
    To struct {
        Email string `json:"email"`
        Name  string `json:"name"`
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
            Size        int64          `json:"size"`
            Encoding    string         `json:"encoding"`
            Compression string         `json:"compression"`
            Checksum    []ChecksumType `json:"checksum"`
            Content     string         `json:"content"`
        } `json:"blocks"`
        Attachments []struct {
            Id          string         `json:"id"`
            MimeType    string         `json:"mimetype"`
            FileName    string         `json:"filename"`
            Size        int64          `json:"size"`
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

//
//
//type MessageProofOfWork struct {
//    Data     int64      `json:"data"`
//}

type ProofOfWorkType struct {
   Bits     int   `json:"bits"`
   Proof    int64 `json:"proof"`
}

//type FromType struct {
//    Email           string              `json:"email"`
//    Name            string              `json:"name"`
//    Organisation    string              `json:"organisation"`
//    ProofOfWork     ProofOfWorkType     `json:"proof_of_work"`
//    PublicKey       string              `json:"public_key"`
//}
//
//type ToType struct {
//    Email           string  `json:"email"`
//    Name            string  `json:"name"`
//}
//
//type BodyHeaderType struct {
//    Size        int64           `json:"size"`
//    Checksum    []ChecksumType  `json:"checksum"`
//}
//
//type MessageHeader struct {
//    From        FromType      `json:"from"`
//    To          ToType        `json:"to"`
//    Body        BodyHeaderType `json:"body"`
//}

type ChecksumType struct {
   Hash    string  `json:"hash"`
   Value   string  `json:"value"`
}

//type Block struct {
//    Size        string              `json:"size"`
//    Encoding    string              `json:"encoding"`
//    Compression string              `json:"compression"`
//    Checksum    []ChecksumType      `json:"checksum"`
//    Body        []byte              `json:"body"`
//}
//
//type Attachment struct {
//    Name            string              `json:"name"`
//    Size            string              `json:"size"`
//    MimeType        string              `json:"mimetype"`
//    Encoding        string              `json:"encoding"`
//    Compression     string              `json:"compression"`
//    Checksum        []ChecksumType      `json:"checksum"`
//    Body            []byte              `json:"body"`
//}
//
//type InputMessageBody struct {
//    ThreadId    string          `json:"thread_id"`
//    Subject     string          `json:"subject"`
//    Labels      []string        `json:"labels"`
//    Blocks      []Block         `json:"blocks"`
//    Attachments []Attachment    `json:"files"`
//}
//
//type Message struct {
//    Body     []byte `json:"message"`
//    From     string `json:"from"`
//    To       string `json:"to"`
//}
//
//type EncryptedMessage struct {
//    Body     []byte `json:"message"`
//    From     string `json:"from"`
//    To       string `json:"to"`
//    Key      []byte `json:"key"`
//    Iv       []byte `json:"iv"`
//}

//func NewMessage(from string, to string, subject string, blocks []InBlock, attachments []InAttachment) *Message {
//    msg := Message{
//        Body: []byte(to),
//        From: from,
//        To: to,
//    }
//
//    return &msg
//}
//
//func NewEncryptedMessage(from string, to string, body []byte, key []byte, iv []byte) *EncryptedMessage {
//    msg := EncryptedMessage{
//        Body: body,
//        From: from,
//        To: to,
//        Key: key,
//        Iv: iv,
//    }
//
//    return &msg
//}
