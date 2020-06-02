package account

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/json"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/config"
    "golang.org/x/crypto/pbkdf2"
    "io"
    "io/ioutil"
    "path"
)

type ProofOfWork struct {
    Bits    int     `json:"bits"`
    Proof   uint64  `json:"proof"`
}

type AccountInfo struct {
    Address         string          `json:"address"`
    Name            string          `json:"name"`
    Organisation    string          `json:"organisation"`
    PrivKey         string          `json:"privKey"`
    PubKey          string          `json:"pubKey"`
    Pow             ProofOfWork     `json:"pow"`
}

type AccountFile struct {
    Address         string  `json:"address"`
    AccountInfo     []byte  `json:"data"`
    Salt            []byte  `json:"salt"`
    Iv              []byte  `json:"iv"`
    Hmac            []byte  `json:"hmac"`
}

const (
    PBKDF_ITERATIONS = 100002
)

func LoadAccount(addr core.Address, password []byte) (*AccountInfo, error) {
    data, err := ioutil.ReadFile(getPath(addr))
    if err != nil {
        return nil, err
    }

    af := &AccountFile{}
    err = json.Unmarshal(data, &af)
    if err != nil {
        return nil, err
    }

    // @TODO: Check HMAC

    derivedAESKey := pbkdf2.Key(password, af.Salt, PBKDF_ITERATIONS, 32, sha256.New)
    aes256, err := aes.NewCipher(derivedAESKey)
    if err != nil {
        return nil, err
    }

    plainText := make([]byte, len(af.AccountInfo))
    ctr := cipher.NewCTR(aes256, af.Iv)
    ctr.XORKeyStream(af.AccountInfo, plainText)

    ai := &AccountInfo{}
    err = json.Unmarshal(plainText, &ai)
    if err != nil {
        return nil, err
    }

    return ai, nil
}

func SaveAccount(addr core.Address, password []byte, acc AccountInfo) (error) {
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
    derivedAESKey := pbkdf2.Key(password, salt, PBKDF_ITERATIONS, 32, sha256.New)
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

    af := &AccountFile{
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

func getPath(addr core.Address) string {
    return path.Clean(path.Join(config.Client.Account.Path, addr.String() + ".account.json"))
}
