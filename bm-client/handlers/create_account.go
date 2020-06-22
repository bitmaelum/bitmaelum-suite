package handlers

import (
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/bm-client/account"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/container"
    "github.com/bitmaelum/bitmaelum-server/core/encrypt"
    "os"
)

func CreateAccount(address, name, organisation, server string, passwd []byte) {
	if !core.IsValidAddress(address) {
		fmt.Printf("Not a valid address")
		fmt.Println("")
		os.Exit(1)
	}

    addr, _ := core.NewAddressFromString(address)
    ks := container.GetResolveService()
    _, err := ks.Resolve(*addr)
    if err == nil {
        fmt.Printf("It seems that this address is already in use. Please specify another address.")
        fmt.Println("")
        os.Exit(1)
    }

    pubKey, privKey, err := encrypt.GenerateKeyPair(encrypt.KeyTypeRSA)
    if err != nil {
        fmt.Print(err)
        fmt.Println("")
        os.Exit(1)
    }

    proof := core.NewProofOfWork(23, addr.Bytes(), 0)

    info := core.AccountInfo{
        Address:      address,
        Name:         name,
        Organisation: organisation,
        PrivKey:      privKey,
        PubKey:       pubKey,
        Pow:          *proof,
        Server:       server,
    }

    account.Vault.Add(info)
    err = account.Vault.Save()
    if err != nil {
        fmt.Printf("Error while saving account into vault: %#v", err)
        fmt.Println("")
        os.Exit(1)
    }

    fmt.Printf("Account created and stored into the vault")
    fmt.Println("")
}
