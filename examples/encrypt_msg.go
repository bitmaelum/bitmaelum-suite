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

const EMAIL = "info@seams-cms.com"

type Message struct {
    Message     string `json:"message"`
    From        string `json:"from"`
    To          string `json:"to"`
    Key         []byte `json:"key"`
    Iv          []byte `json:"iv"`
}

func main() {
    // Create a message, encrypt it with somebody's key and save to disk
    var msg = Message{
        Message: "Hi there.. this is an encrypted message that you should not be able to decrypt",
        From: EMAIL,
        To: "joshua@noxlogic.nl",
    }

    privKey, err := LoadPrivateKey(EMAIL)
    if err != nil {
        panic(err)
    }

    err = encrypt(&msg, privKey)
    if err != nil {
        panic(err)
    }

    jsonData, err := json.MarshalIndent(msg, "", " ")
    if err != nil {
        panic(err)
    }
    err = ioutil.WriteFile("message.enc", jsonData, 0644)
    if err != nil {
        panic(err)
    }
}

func encrypt(msg *Message, privKey *rsa.PrivateKey) error {
    //
    // First, encrypt the data with symmetrical encryption (AES) and a randomized key and IV
    //

    msg.Iv = make([]byte, 16)
    _, err := rand.Read(msg.Iv)
    if err != nil {
        return err
    }

    aesKey := make([]byte, 32)
    _, err = rand.Read(aesKey)
    if err != nil {
        return err
    }
    block, err := aes.NewCipher([]byte(aesKey))
    if err != nil {
        return err
    }
    plaintext := []byte(msg.Message)
    cfb := cipher.NewCFBEncrypter(block, msg.Iv)
    ciphertext := make([]byte, len(plaintext))
    cfb.XORKeyStream(ciphertext, plaintext)
    msg.Message = base64.StdEncoding.EncodeToString(ciphertext)


    //
    // Next, encrypt the symmetrical key with PUB/PRIV key system (asymmetrical)
    //

    // Load the public key from the recipient. This will allow only the PRIV
    // key holder (recipient) to decode the data.
    publicKey, err := utils.LoadPubKey(fmt.Sprintf("%s.pub", msg.To))
    if err != nil {
        return err
    }

    msg.Key, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
    if err != nil {
        return err
    }

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
