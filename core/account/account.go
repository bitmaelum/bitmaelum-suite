package account

import (
    "encoding/json"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/config"
    "io/ioutil"
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

func LoadAccount(addr core.Address) (*AccountInfo, error) {
    path := config.Client.Account.Path + addr.String() + ".account.json"

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    ai := &AccountInfo{}
    err = json.Unmarshal(data, &ai)
    if err != nil {
        return nil, err
    }

    return ai, nil
}