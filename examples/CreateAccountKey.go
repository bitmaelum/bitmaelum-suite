package examples

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/rsa"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/jaytaph/mailv2/utils"
    "io/ioutil"
)

//const EMAIL = "info@seams-cms.com"
const EMAIL = "joshua@noxlogic.nl"

type Message struct {
    Message     string `json:"message"`
    From        string `json:"from"`
    To          string `json:"to"`
    Key         []byte `json:"key"`
    Iv          []byte `json:"iv"`
}

func main() {
    msg, err := loadMessage("message.enc")
    if err != nil {
        panic(err)
    }

    privKey, err := LoadPrivateKey(EMAIL)
    if err != nil {
        panic(err)
    }

    err = decrypt(msg, privKey)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s", msg.Message)
}

func loadMessage(path string) (*Message, error) {
    jsonData, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    msg := Message{}
    err = json.Unmarshal(jsonData, &msg)
    if err != nil {
        return nil, err
    }

    return &msg, nil
}

func decrypt(msg *Message, privKey *rsa.PrivateKey) error {
    //
    // Decrypt AES key
    //
    aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privKey, []byte(msg.Key))
    if err != nil {
        return err
    }

    //
    // Use AES key and IV to decode message
    //
    block, err := aes.NewCipher([]byte(aesKey))
    if err != nil {
        return err
    }
    ciphertext, err := base64.StdEncoding.DecodeString(msg.Message)
    if err != nil {
        return err
    }
    cfb := cipher.NewCFBDecrypter(block, msg.Iv)
    plaintext := make([]byte, len(ciphertext))
    cfb.XORKeyStream(plaintext, ciphertext)
    msg.Message = string(plaintext)

    return nil
}

func LoadPrivateKey(email string) (*rsa.PrivateKey, error) {
    privateKey, err := utils.LoadPrivKey(fmt.Sprintf("%s.key.pem", email))
    if err != nil {
        privateKey, err := utils.CreateNewKeyPair(4096)
        if err != nil {
            panic(err)
        }

        err = utils.SavePrivKey(fmt.Sprintf("%s.key.pem", email), privateKey)
        if err != nil {
            panic(err)
        }
    }

    return privateKey, nil
}

