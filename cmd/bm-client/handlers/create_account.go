package handlers

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/pkg/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"os"
)

// CreateAccount creates a new account locally in the vault, stores it on the mailserver and pushes the public key to the resolver
func CreateAccount(vault *vault.Vault, bmAddr, name, organisation, server, token string) {

	fmt.Printf("* Verifying if address is correct: ")
	addr, err := address.New(bmAddr)
	if err != nil {
		fmt.Printf("not a valid address")
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("ok\n")

	fmt.Printf("* Checking if address is already known in the resolver service: ")
	ks := container.GetResolveService()
	_, err = ks.Resolve(addr.Hash())
	if err == nil {
		fmt.Printf("it seems that this address is already in use. Please specify another address.")
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("not found. This is a good thing.\n")

	fmt.Printf("* Checking if the account is already present in the vault: ")
	if vault.HasAccount(*addr) {
		fmt.Printf("account already present in the vault")
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("not found. This is a good thing.\n")

	fmt.Printf("* Generating your secret key to send and read mail: ")
	privKey, pubKey, err := encrypt.GenerateKeyPair(bmcrypto.KeyTypeRSA)
	if err != nil {
		fmt.Print(err)
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done.\n")

	fmt.Printf("* Doing some work to let people know this is not a fake account: ")
	proof := pow.New(config.Client.Accounts.ProofOfWork, addr.Hash().String(), 0)
	proof.Work()
	fmt.Printf("done.\n")

	fmt.Printf("* Adding your new account into the vault: ")
	info := pkg.Info{
		Address: bmAddr,
		Name:    name,
		PrivKey: *privKey,
		PubKey:  *pubKey,
		Pow:     *proof,
		Server:  server,
	}

	vault.AddAccount(info)
	err = vault.Save()
	if err != nil {
		fmt.Printf("error while saving account into vault: %#v", err)
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done\n")

	fmt.Printf("* Sending your account information to the server: ")
	client, err := api.NewAuthenticated(&info, api.ClientOpts{
		Host:          info.Server,
		AllowInsecure: config.Client.Server.AllowInsecure,
	})
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		vault.RemoveAccount(*addr)
		_ = vault.Save()

		fmt.Printf("cannot initialize API")
		fmt.Println("")
		os.Exit(1)
	}

	err = client.CreateAccount(info, token)
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		vault.RemoveAccount(*addr)
		_ = vault.Save()

		fmt.Printf("error from API while trying to create account: " + err.Error())
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done\n")

	fmt.Printf("* Making your account known to the outside world: ")
	err = ks.UploadInfo(info, info.Server)
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		vault.RemoveAccount(*addr)
		_ = vault.Save()

		fmt.Printf("error while uploading account to the resolver: " + err.Error())
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done\n")

	fmt.Printf("\n")
	fmt.Printf("* All done")
}
