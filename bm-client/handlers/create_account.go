package handlers

import (
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/bm-client/account"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/api"
    "github.com/bitmaelum/bitmaelum-server/core/container"
    "github.com/bitmaelum/bitmaelum-server/core/encrypt"
    "os"
)

func CreateAccount(address, name, organisation, server, token string, passwd []byte) {
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

    if account.Vault.HasAccount(*addr) {
        fmt.Printf("Account already present in the vault")
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


    client, err := api.CreateNewClient(&info)
    if err != nil {
        // Remove account from the local vault as well, as we could not store on the server
        account.Vault.Remove(*addr)
        _ = account.Vault.Save()

        fmt.Printf("Cannot initialize API")
        fmt.Println("")
        os.Exit(1)
    }

    err = client.CreateAccount(info, token)
    if err != nil {
        // Remove account from the local vault as well, as we could not store on the server
        account.Vault.Remove(*addr)
        _ = account.Vault.Save()

        fmt.Printf("Error from API while trying to create account: " + err.Error())
        fmt.Println("")
        os.Exit(1)
    }

    fmt.Printf("Account created and stored into the vault")
    fmt.Println("")
}
