package handlers

import (
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/invite"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

// CreateAccount creates a new account locally in the vault, stores it on the mail server and pushes the public key to the resolver
func CreateAccount(vault *vault.Vault, bmAddr, name, token string) {
	fmt.Printf("* Verifying if address is correct: ")
	addr, err := address.NewAddress(bmAddr)
	if err != nil {
		fmt.Printf("not a valid address")
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("ok\n")

	if addr.HasOrganisationPart() {
		fmt.Printf("* You are creating an organisation address.")
	}

	fmt.Printf("* Checking if address is already known in the resolver service: ")
	ks := container.GetResolveService()
	_, err = ks.ResolveAddress(addr.Hash())
	if err == nil {
		fmt.Printf("\n  X it seems that this address is already in use. Please specify another address.")
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("not found. This is a good thing.\n")

	// Check token
	fmt.Printf("* Checking token format and extracting data: ")
	it, err := invite.ParseInviteToken(token)
	if err != nil {
		fmt.Printf("\n  X it seems that this token is invalid")
		fmt.Println("")
		os.Exit(1)
	}
	// Check address matches the one in the token
	if it.Address.String() != addr.Hash().String() {
		fmt.Printf("\n  X this token is not for %s", addr.String())
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Printf("token is valid.\n")

	fmt.Printf("* Checking if the account is already present in the vault: ")
	var info *internal.AccountInfo
	var seedToShow string
	if vault.HasAccount(*addr) {
		fmt.Printf("\n  X account already present in the vault. Strange, but let's continue...\n")
		info, err = vault.GetAccountInfo(*addr)
		if err != nil {
			fmt.Print(err)
			fmt.Println("")
			os.Exit(1)
		}
	} else {
		fmt.Printf("not found. This is a good thing.\n")

		fmt.Printf("* Generating your secret key to send and read mail: ")
		seed, privKey, pubKey, err := bmcrypto.GenerateKeypairWithSeed()
		if err != nil {
			fmt.Print(err)
			fmt.Println("")
			os.Exit(1)
		}
		seedToShow = seed
		fmt.Printf("done.\n")

		fmt.Printf("* Doing some work to let people know this is not a fake account: ")
		proof := pow.NewWithoutProof(config.Client.Accounts.ProofOfWork, addr.Hash().String())
		proof.WorkMulticore()
		fmt.Printf("done.\n")

		fmt.Printf("* Adding your new account into the vault: ")
		info = &internal.AccountInfo{
			Address:   bmAddr,
			Name:      name,
			PrivKey:   *privKey,
			PubKey:    *pubKey,
			Pow:       *proof,
			RoutingID: it.RoutingID, // Fetch from token
		}

		vault.AddAccount(*info)
		err = vault.WriteToDisk()
		if err != nil {
			fmt.Printf("\n  X error while saving account into vault: %#v", err)
			fmt.Println("")
			os.Exit(1)
		}
		fmt.Printf("done\n")
	}

	// Fetch routing info
	routingInfo, err := ks.ResolveRouting(it.RoutingID)
	if err != nil {
		fmt.Printf("\n  X routing information is not found: %#v", err)
		fmt.Println("")
		os.Exit(1)
	}

	fmt.Printf("* Sending your account information to the server: ")
	client, err := api.NewAuthenticated(info, api.ClientOpts{
		Host:          routingInfo.Routing,
		AllowInsecure: config.Client.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
	})
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		vault.RemoveAccount(*addr)
		_ = vault.WriteToDisk()

		fmt.Printf("cannot initialize API")
		fmt.Println("")
		os.Exit(1)
	}

	err = client.CreateAccount(*info, token)
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		vault.RemoveAccount(*addr)
		_ = vault.WriteToDisk()

		fmt.Printf("\n  X error from API while trying to create account: " + err.Error())
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done\n")

	fmt.Printf("* Making your account known to the outside world: ")
	err = ks.UploadAddressInfo(*info)
	if err != nil {
		// We can't remove the account from the vault as we have created it on the mail-server

		fmt.Printf("\n  X error while uploading account to the resolver: " + err.Error())
		fmt.Printf("\n  X Please try again with:\n   bm-client push-account -a '%s'\n", addr.String())
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done\n")

	fmt.Printf("\n")
	fmt.Printf("* All done")

	if len(seedToShow) > 0 {
		fmt.Print(`
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your account. 
If, for any reason, you lose this key, you will need to use the following 
words in order to recreate the key:

`)
		fmt.Print(internal.WordWrap(seedToShow, 78))
		fmt.Print(`

Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.

WITHOUT THESE WORDS, ALL ACCESS TO YOUR ACCOUNT IS LOST!
`)
	}
}
