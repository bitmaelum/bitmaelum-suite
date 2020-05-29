package account

import (
    "crypto/x509"
    "encoding/json"
    "encoding/pem"
    "errors"
    "fmt"
    "github.com/jaytaph/mailv2/core/utils"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "path"
)

type ProofOfWork struct {
    Bits    int     `json:"bits"`
    Proof   uint64  `json:"proof"`
}

type Account struct {
    Address         string          `json:"address"`
    Name            string          `json:"name"`
    Organisation    string          `json:"organisation"`
    PrivKey         string          `json:"privKey"`
    PubKey          string          `json:"pubKey"`
    Pow             ProofOfWork     `json:"pow"`
}

type AccountInfo struct {
    Accounts []Account
}

const ACCOUNT_FILE = "accounts.json"

func getAccountPath() string {

    /*
    homedir, err := os.UserHomeDir()
    if err != nil {
       homedir = "."
    }

    accountPath := fmt.Sprintf("%s/.mv2/%s", homedir, ACCOUNT_FILE)
    */

    return fmt.Sprintf("./%s", ACCOUNT_FILE)
}

func (a *AccountInfo) SaveAccount() {
    accountPath := getAccountPath()

    if _, err := os.Stat(accountPath); os.IsNotExist(err) {
        // Create path
        dir := path.Dir(accountPath)
        err := os.MkdirAll(dir, 700)
        if err != nil {
            logrus.Error(err)
            os.Exit(128)
        }
    }

    // Marshal and save account data
    data, err := json.MarshalIndent(a.Accounts, "", " ")
    if err != nil {
        logrus.Error(err)
        os.Exit(128)
    }

    err = ioutil.WriteFile(accountPath, data, 0600)
    if err != nil {
        logrus.Error(err)
        os.Exit(128)
    }
}

func LoadAccountConfig() *AccountInfo {
    accountPath := getAccountPath()

    if _, err := os.Stat(accountPath); os.IsNotExist(err) {
        return &AccountInfo{}
    }

    data, err := ioutil.ReadFile(accountPath)
    if err != nil {
        panic(err)
    }

    accountInfo := AccountInfo{}
    err = json.Unmarshal(data, &accountInfo.Accounts)
    if err != nil {
        logrus.Errorf("error while parsing account file: ", err)
    }

    return &accountInfo
}


func (a *AccountInfo) Has(address string) bool {
    for idx := range a.Accounts {
        if a.Accounts[idx].Address == address {
            return true
        }
    }
    return false
}

func (a *AccountInfo) Get(address string) (*Account, error) {
    for idx := range a.Accounts {
        if a.Accounts[idx].Address == address {
            return &a.Accounts[idx], nil
        }
    }

    return nil, errors.New("account not found")
}

func (a *AccountInfo) GenerateAccount(address string, name string) (*Account, error) {
    logrus.Trace("generating new keypair")
    privateKey, err := utils.CreateNewKeyPair(4096)
    if err != nil {
        return nil, err
    }

    logrus.Trace("calculating for proof-of-work")
    pow := ProofOfWork{
        Bits:  20,
        Proof: utils.ProofOfWork(20, []byte(address)),
    }

    privPem := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })
    pubPem := pem.EncodeToMemory(&pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
    })

    account := Account{
        Address: address,
        Name:    name,
        PrivKey: string(privPem),
        PubKey:  string(pubPem),
        Pow:     pow,
    }

    logrus.Debug("saving account")
    a.Accounts = append(a.Accounts, account)
    a.SaveAccount()

    return &account, nil
}
