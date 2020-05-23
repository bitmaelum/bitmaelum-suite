package message

type ProofOfWorkType struct {
    Bits     int64 `json:"bits"`
    Nonce    int64 `json:"nonce"`
}

type FromType struct {
    Email           string              `json:"email"`
    Name            string              `json:"name"`
    Organisation    string              `json:"organisation"`
    ProofOfWork     ProofOfWorkType     `json:"proof_of_work"`
    PublicKey       string              `json:"public_key"`
}

type ToType struct {
    Email           string  `json:"email"`
    Name            string  `json:"name"`
}

type BodyHeaderType struct {
    Size        int64           `json:"size"`
    Checksum    []ChecksumType  `json:"checksum"`
}

type MessageHeader struct {
    From        FromType      `json:"from"`
    To          ToType        `json:"to"`
    Body        BodyHeaderType `json:"body"`
}

type ChecksumType struct {
    Hash    string  `json:"hash"`
    Value   string  `json:"value"`
}

type Block struct {
    Size        string              `json:"size"`
    Encoding    string              `json:"encoding"`
    Compression string              `json:"compression"`
    Checksum    []ChecksumType      `json:"checksum"`
    Body        []byte              `json:"body"`
}

type File struct {
    Name            string              `json:"name"`
    Size            string              `json:"size"`
    MimeType        string              `json:"mimetype"`
    Encoding        string              `json:"encoding"`
    Compression     string              `json:"compression"`
    Checksum        []ChecksumType      `json:"checksum"`
    Body            []byte              `json:"body"`
}

type InputMessageBody struct {
    ThreadId    string      `json:"thread_id"`
    Subject     string      `json:"subject"`
    Labels      []string    `json:"labels"`
    Blocks      []Block     `json:"blocks"`
    Files       []File      `json:"files"`
}

type Message struct {
    Body     []byte `json:"message"`
    From     string `json:"from"`
    To       string `json:"to"`
}

type EncryptedMessage struct {
    Body     []byte `json:"message"`
    From     string `json:"from"`
    To       string `json:"to"`
    Key      []byte `json:"key"`
    Iv       []byte `json:"iv"`
}

func NewMessage(from string, to string, body []byte) *Message {
    msg := Message{
        Body: body,
        From: from,
        To: to,
    }

    return &msg
}

func NewEncryptedMessage(from string, to string, body []byte, key []byte, iv []byte) *EncryptedMessage {
    msg := EncryptedMessage{
        Body: body,
        From: from,
        To: to,
        Key: key,
        Iv: iv,
    }

    return &msg
}
