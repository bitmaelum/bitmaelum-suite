package handlers

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/vault"
	"github.com/bitmaelum/bitmaelum-suite/core/api"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/core/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"os"
)

// CreateAccount creates a new account locally in the vault, stores it on the mailserver and pushes the public key to the resolver
func CreateAccount(vault *vault.Vault, bmAddr, name, organisation, server, token string) {
	addr, err := address.New(bmAddr)
	if err != nil {
		fmt.Printf("Not a valid address")
		fmt.Println("")
		os.Exit(1)
	}

	ks := container.GetResolveService()
	_, err = ks.Resolve(addr.Hash())
	if err == nil {
		fmt.Printf("It seems that this address is already in use. Please specify another address.")
		fmt.Println("")
		os.Exit(1)
	}

	if vault.HasAccount(*addr) {
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

	proof := pow.New(config.Client.Accounts.ProofOfWork, []byte(addr.Hash()), 0)
	proof.Work()

	info := account.Info{
		Address:      bmAddr,
		Name:         name,
		Organisation: organisation,
		PrivKey:      privKey,
		PubKey:       pubKey,
		Pow:          *proof,
		Server:       server,
	}

	vault.AddAccount(info)
	err = vault.Save()
	if err != nil {
		fmt.Printf("Error while saving account into vault: %#v", err)
		fmt.Println("")
		os.Exit(1)
	}

	client, err := api.CreateNewClient(&info)
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		vault.RemoveAccount(*addr)
		_ = vault.Save()

		fmt.Printf("Cannot initialize API")
		fmt.Println("")
		os.Exit(1)
	}

	err = client.CreateAccount(info, token)
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		vault.RemoveAccount(*addr)
		_ = vault.Save()

		fmt.Printf("Error from API while trying to create account: " + err.Error())
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Printf("Account created and stored into the vault")
	fmt.Println("")
}
