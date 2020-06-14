package account

import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/json"
    "errors"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/api"
    "github.com/bitmaelum/bitmaelum-server/core/config"
    "github.com/mitchellh/go-homedir"
    "github.com/sirupsen/logrus"
    "golang.org/x/crypto/pbkdf2"
    "io"
    "io/ioutil"
    "path"
)

type EncryptedAccountInfo struct {
    Address         string  `json:"address"`
    AccountInfo     []byte  `json:"data"`
    Salt            []byte  `json:"salt"`
    Iv              []byte  `json:"iv"`
    Hmac            []byte  `json:"hmac"`
}

const (
    PbkdfIterations = 100002
)

// Load local account information
func LoadAccount(addr core.Address, password []byte) (*core.AccountInfo, error) {
    data, err := ioutil.ReadFile(getPath(addr))
    if err != nil {
        return nil, err
    }

    var plainText []byte

    af := &EncryptedAccountInfo{}
    err = json.Unmarshal(data, &af)
    if err != nil || af.AccountInfo == nil {
        logrus.Warning("account file is unprotected with a password.")

        plainText = []byte(data)
    } else {
        // Check HMAC
        hash := hmac.New(sha256.New, password)
        hash.Write(af.AccountInfo)
        if bytes.Compare(hash.Sum(nil), af.Hmac) != 0 {
            return nil, errors.New("HMAC incorrect")
        }

        derivedAESKey := pbkdf2.Key(password, af.Salt, PbkdfIterations, 32, sha256.New)
        aes256, err := aes.NewCipher(derivedAESKey)
        if err != nil {
            return nil, err
        }

        plainText = make([]byte, len(af.AccountInfo))
        ctr := cipher.NewCTR(aes256, af.Iv)
        ctr.XORKeyStream(plainText, af.AccountInfo)
    }

    ai := &core.AccountInfo{}
    err = json.Unmarshal(plainText, &ai)
    if err != nil {
        return nil, err
    }

    return ai, nil
}

// Create local account
func CreateLocalAccount(addr core.Address, password []byte, acc core.AccountInfo) (error) {
    // Generate JSON structure that we will encrypt
    plainText, err := json.MarshalIndent(&acc, "", "  ")
    if err != nil {
        return err
    }

    // Generate 64 byte salt
    salt := make([]byte, 64)
    _, err = io.ReadFull(rand.Reader, salt)
    if err != nil {
        return err
    }

    // Generate key based on password
    derivedAESKey := pbkdf2.Key(password, salt, PbkdfIterations, 32, sha256.New)
    aes256, err := aes.NewCipher(derivedAESKey)
    if err != nil {
        return err
    }

    // Generate 32 byte IV
    iv := make([]byte, aes.BlockSize)
    _, err = io.ReadFull(rand.Reader, iv)
    if err != nil {
        return err
    }

    // Encrypt the data
    cipherText := make([]byte, len(plainText))
    ctr := cipher.NewCTR(aes256, iv)
    ctr.XORKeyStream(cipherText, plainText)


    // Generate HMAC
    hash := hmac.New(sha256.New, password)
    hash.Write(cipherText)

    af := &EncryptedAccountInfo{
        Address:     addr.String(),
        AccountInfo: cipherText,
        Salt:        salt,
        Iv:          iv,
        Hmac:        hash.Sum(nil),
    }
    data, err := json.MarshalIndent(af, "", "  ")
    if err != nil {
        return err
    }

    // And write to file
    err = ioutil.WriteFile(getPath(addr), data, 0600)
    if err != nil {
        return err
    }

    return nil
}

// Create remote account
func CreateRemoteAccount(acc *core.AccountInfo, token string) (error) {
    client, err := api.CreateNewClient(acc)
    if err != nil {
        return err
    }

    return client.CreateAccount(core.StringToHash(acc.Address), token)
}

// Generate path to account file for the given address
func getPath(addr core.Address) string {
    p := path.Join(config.Client.Accounts.Path, addr.String() + ".account.json")
    p, _ = homedir.Expand(p)
    return p
}
