package client

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
    "github.com/bitmaelum/bitmaelum-server/core/password"
    "github.com/mitchellh/go-homedir"
    "github.com/opentracing/opentracing-go/log"
    "github.com/sirupsen/logrus"
    "golang.org/x/crypto/pbkdf2"
    "io"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "strings"
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

// @TODO: We probably want to use scrypt instead of PBKDF

// LoadAccount loads either an encrypted or unencrypted account info file
func LoadAccount(addr core.Address, pwd []byte) (*core.AccountInfo, error) {
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
        hash := hmac.New(sha256.New, pwd)
        hash.Write(af.AccountInfo)
        if bytes.Compare(hash.Sum(nil), af.Hmac) != 0 {
            return nil, errors.New("HMAC incorrect")
        }

        derivedAESKey := pbkdf2.Key(pwd, af.Salt, PbkdfIterations, 32, sha256.New)
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
func CreateLocalAccount(addr core.Address, pwd []byte, account core.AccountInfo) (error) {
    plainText, err := json.MarshalIndent(&account, "", "  ")
    if err != nil {
        return err
    }

    return safeWrite(getPath(addr), plainText)
}

// Create remote account
func CreateRemoteAccount(acc *core.AccountInfo, token string) (error) {
    client, err := api.CreateNewClient(acc)
    if err != nil {
        return err
    }

    return client.CreateAccount(core.StringToHash(acc.Address), token)
}

func GetAccounts() ([]core.AccountInfo, error) {
    accounts := []core.AccountInfo{}

    // Read whole accounts directory
    p, _ := homedir.Expand(config.Client.Accounts.Path)
    files, err := ioutil.ReadDir(p)
    if err != nil {
        return nil, err
    }

    // Iterate each file
    for _, f := range files {
        // Check if the file is an account file
        matched, err := filepath.Match("*.account.json", f.Name())
        if !matched || err != nil {
            continue
        }

        // Fetch address from account filename
        address := strings.Replace(f.Name(), ".account.json", "", 1)
        addr, err := core.NewAddressFromString(address)
        if err != nil {
            log.Error(err)
            continue;
        }

        // Load account
        pwd, err := password.FetchPassword(addr)
        if err != nil {
            log.Error(err)
            continue
        }
        acc, err := LoadAccount(*addr, pwd)
        if err != nil {
            log.Error(err)
            continue
        }

        err = password.StorePassword(addr, pwd)
        if err != nil {
            log.Error(err)
        }

        accounts = append(accounts, *acc)
    }

    return accounts, nil
}

// Locks an account on disk with given pwd
func LockAccount(addr core.Address, pwd []byte) error {
    // Load account from disk. Assumes unlocked account.
    account, err := LoadAccount(addr, []byte{})
    if err != nil {
        return err
    }

    // Encrypt account
    encryptedAccount, err := encryptAccount(*account, pwd)
    if err != nil {
        return err
    }

    // Write file to disk
    data, err := json.MarshalIndent(encryptedAccount, "", "  ")
    if err != nil {
        return err
    }

    return safeWrite(getPath(addr), data)
}

// Returns true when the given account is locked on disk (without unlocking it)
func IsLocked(addr core.Address) bool {
    data, err := ioutil.ReadFile(getPath(addr))
    if err != nil {
        return false
    }

    af := &EncryptedAccountInfo{}
    err = json.Unmarshal(data, &af)

    return err != nil
}

// Locks the given account to disk with the given password
func UnlockAccount(addr core.Address, pwd []byte) error {
    account, err := LoadAccount(addr, pwd)
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(&account, "", "  ")
    if err != nil {
        return nil
    }

    return safeWrite(getPath(addr), data)
}

// Generate path to account file for the given address
func getPath(addr core.Address) string {
    p := path.Join(config.Client.Accounts.Path, addr.String() + ".account.json")
    p, _ = homedir.Expand(p)
    return p
}

// Converts accountInfo struct into EncryptedAccountInfo struct
func encryptAccount(account core.AccountInfo, pwd []byte) (*EncryptedAccountInfo, error){
    // Generate 64 byte salt
    salt := make([]byte, 64)
    _, err := io.ReadFull(rand.Reader, salt)
    if err != nil {
        return nil, err
    }

    // Generate key based on password
    derivedAESKey := pbkdf2.Key(pwd, salt, PbkdfIterations, 32, sha256.New)
    aes256, err := aes.NewCipher(derivedAESKey)
    if err != nil {
        return nil, err
    }

    // Generate 32 byte IV
    iv := make([]byte, aes.BlockSize)
    _, err = io.ReadFull(rand.Reader, iv)
    if err != nil {
        return nil, err
    }

    // Encrypt the data
    plainText, err := json.MarshalIndent(&account, "", "  ")
    if err != nil {
        return nil, err
    }

    cipherText := make([]byte, len(plainText))
    ctr := cipher.NewCTR(aes256, iv)
    ctr.XORKeyStream(cipherText, plainText)


    // Generate HMAC
    hash := hmac.New(sha256.New, pwd)
    hash.Write(cipherText)

    return &EncryptedAccountInfo{
        Address:     account.Address,
        AccountInfo: cipherText,
        Salt:        salt,
        Iv:          iv,
        Hmac:        hash.Sum(nil),
    }, nil
}

// Writes data by safely writing to a temp file first
func safeWrite(path string, data []byte) error {
    err := ioutil.WriteFile(path + ".tmp", data, 0600)
    if err != nil {
        return err
    }

    err = os.Rename(path + ".tmp", path)
    return err
}