package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/pkg/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
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
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	fmt.Println(internal.GetASCIILogo())

	// Unlock vault
	accountVault, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	addr, err := address.New(opts.Addr)
	if err != nil {
		panic(err)
	}

	info, err := accountVault.GetAccountInfo(*addr)
	if err != nil {
		panic(err)
	}
	hash, err := address.NewHash(info.Address)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Reading message for user %s (%s) (%s)\n", info.Name, info.Address, hash.String())

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

	privKey, err := encrypt.PEMToPrivKey([]byte(info.PrivKey))
	if err != nil {
		panic(err)
	}

	decryptedKey, err := encrypt.Decrypt(privKey, header.Catalog.EncryptedKey)
	if err != nil {
		panic(err)
	}

	catalog, err := encrypt.CatalogDecrypt(decryptedKey, data)
	if err != nil {
		panic(err)
	}

	s, _ := json.MarshalIndent(catalog, "", "  ")
	fmt.Printf("%s", s)

	// Read, decrypt and display blocks
	for _, block := range catalog.Blocks {
		f, err := os.Open(opts.Path + "/" + block.ID)
		if err != nil {
			panic(err)
		}

		r, err := message.GetAesDecryptorReader(block.IV, block.Key, f)
		if err != nil {
			f.Close()
			panic(err)
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			f.Close()
			panic(err)
		}

		fmt.Printf("\n----- START BLOCK (%s) %s --------\n", block.Type, block.ID)
		fmt.Printf("%s", content)
		fmt.Printf("\n----- END BLOCK %s --------\n", block.ID)
	}
}
