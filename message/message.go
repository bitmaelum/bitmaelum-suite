package message

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
