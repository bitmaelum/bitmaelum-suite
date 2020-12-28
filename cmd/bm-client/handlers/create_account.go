// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package handlers

import (
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	bminternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/signature"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

func verifyAddress(bmAddr string) *address.Address {
	fmt.Printf("* Verifying if address is correct: ")
	addr, err := address.NewAddress(bmAddr)
	if err != nil {
		fmt.Println("not a valid address")
		os.Exit(1)
	}
	fmt.Println("ok")

	if addr.HasOrganisationPart() {
		fmt.Println("* You are creating an organisation address.")
	}

	return addr
}

func checkAddressInResolver(addr address.Address) *resolver.Service {
	fmt.Printf("* Checking if address is already known in the resolver service: ")

	ks := container.Instance.GetResolveService()
	_, err := ks.ResolveAddress(addr.Hash())

	if err == nil {
		fmt.Println("")
		fmt.Println("  X it seems that this address is already in use. Please specify another address.")
		os.Exit(1)
	}

	fmt.Println("not found. This is a good thing.")

	return ks
}

func checkToken(token string, addr address.Address) *signature.Token {
	fmt.Printf("* Checking token format and extracting data: ")

	it, err := signature.ParseInviteToken(token)
	if err != nil {
		fmt.Println("\n  X it seems that this token is invalid")
		os.Exit(1)
	}

	// Check address matches the one in the token
	if it.AddrHash.String() != addr.Hash().String() {
		fmt.Println("\n  X this token is not for", addr.String())
		os.Exit(1)
	}

	fmt.Printf("token is valid.\n")

	return it
}

func checkAccountInVault(vault *vault.Vault, addr address.Address) *vault.AccountInfo {
	fmt.Printf("* Checking if the account is already present in the vault: ")

	if vault.HasAccount(addr) {
		fmt.Printf("\n  X account already present in the vault. Strange, but let's continue...\n")
		info, err := vault.GetAccountInfo(addr)
		if err != nil {
			fmt.Print(err)
			fmt.Println("")
			os.Exit(1)
		}
		return info
	}
	fmt.Printf("not found. This is a good thing.\n")

	return nil
}

// CreateAccount creates a new account locally in the vault, stores it on the mail server and pushes the public key to the resolver
func CreateAccount(v *vault.Vault, bmAddr, name, token string, kt bmcrypto.KeyType) {
	var mnemonicToShow string

	fmt.Println("")

	addr := verifyAddress(bmAddr)
	ks := checkAddressInResolver(*addr)
	it := checkToken(token, *addr)
	info := checkAccountInVault(v, *addr)

	res := container.Instance.GetResolveService()
	resolverCfg := res.GetConfig()

	if info == nil {
		fmt.Printf("* Generating your secret key to send and read mail: ")

		mnemonic, privKey, pubKey, err := bminternal.GenerateKeypairWithMnemonic(kt)

		if err != nil {
			fmt.Print(err)
			fmt.Println("")
			os.Exit(1)
		}
		mnemonicToShow = mnemonic
		fmt.Printf("done.\n")

		fmt.Printf("* Doing some work to let people know this is not a fake account, this might take a while: ")
		proof := pow.NewWithoutProof(resolverCfg.ProofOfWork.Address, addr.Hash().String())
		proof.WorkMulticore()
		fmt.Printf("done.\n")

		fmt.Printf("* Adding your new account into the vault: ")
		info = &vault.AccountInfo{
			Address:   addr,
			Name:      name,
			PrivKey:   *privKey,
			PubKey:    *pubKey,
			Pow:       proof,
			RoutingID: it.RoutingID, // Fetch from token
		}

		v.AddAccount(*info)
		err = v.WriteToDisk()
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
	client, err := api.NewAuthenticated(*info.Address, &info.PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		v.RemoveAccount(*addr)
		_ = v.WriteToDisk()

		fmt.Printf("cannot initialize API")
		fmt.Println("")
		os.Exit(1)
	}

	err = client.CreateAccount(*info, token)
	if err != nil {
		// Remove account from the local vault as well, as we could not store on the server
		v.RemoveAccount(*addr)
		_ = v.WriteToDisk()

		fmt.Printf("\n  X error from API while trying to create account: " + err.Error())
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done\n")

	fmt.Printf("* Making your account known to the outside world: ")
	err = ks.UploadAddressInfo(*info, "")
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

	if len(mnemonicToShow) > 0 {
		fmt.Print(`
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your account. 
If, for any reason, you lose this key, you will need to use the following 
words in order to recreate the key:

`)
		fmt.Print(bminternal.WordWrap(mnemonicToShow, 78))
		fmt.Print(`

Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.

WITHOUT THESE WORDS, ALL ACCESS TO YOUR ACCOUNT IS LOST!
`)
	}
}
