package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/bm-client/account"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/bitmaelum/bitmaelum-server/core/encrypt"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"io/ioutil"
	"os"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Configuration file"`
	Addr     string `short:"a" long:"address" description:"account"`
	Password string `short:"p" long:"password" description:"Password to decrypt your account"`
	Path     string `short:"" long:"path" description:"path to message"`
}

var opts options

func main() {
	core.ParseOptions(&opts)
	core.LoadClientConfig(opts.Config)

	fmt.Println(core.GetAsciiLogo())

	// Unlock vault
	err := account.UnlockVault(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	addr, err := core.NewAddressFromString(opts.Addr)
	if err != nil {
		panic(err)
	}

	ai, err := account.Vault.FindAccount(*addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reading message for user %s (%s) (%s)\n", ai.Name, ai.Address, core.StringToHash(ai.Address))

	data, err := ioutil.ReadFile(opts.Path + "/header.json")
	if err != nil {
		panic(err)
	}

	header := &message.Header{}
	err = json.Unmarshal(data, header)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Reading message from: %s\n", header.From.Addr)
	fmt.Printf("Reading message to: %s\n", header.To.Addr)

	fmt.Printf("Decrypting catalog")
	data, err = ioutil.ReadFile(opts.Path + "/catalog")
	if err != nil {
		panic(err)
	}

	privKey, err := encrypt.PEMToPrivKey([]byte(ai.PrivKey))
	if err != nil {
		panic(err)
	}

	decryptedKey, err := encrypt.Decrypt(privKey, header.Catalog.EncryptedKey)
	if err != nil {
		panic(err)
	}

	catalog, err := encrypt.DecryptCatalog(decryptedKey, data)
	if err != nil {
		panic(err)
	}

	s, _ := json.MarshalIndent(catalog, "", "  ")
	fmt.Printf("%s", s)

	// Read, decrypt and display blocks
	for _, block := range catalog.Blocks {
		f, err := os.Open(opts.Path + "/" + block.Id)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		r, err := message.GetAesDecryptorReader(block.Iv, block.Key, f)
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			panic(err)
		}

		fmt.Printf("\n----- START BLOCK (%s) %s --------\n", block.Type, block.Id)
		fmt.Printf("%s", content)
		fmt.Printf("\n----- END BLOCK %s --------\n", block.Id)
	}
}
