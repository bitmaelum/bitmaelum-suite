package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/vault"
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"log"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
	Addr     string `short:"a" long:"address" description:"account"`
	Password string `short:"p" long:"password" description:"Password to decrypt your account"`
}

var opts options

func main() {
	core.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	// Convert strings into addresses
	fromAddr, err := address.New(opts.Addr)
	if err != nil {
		log.Fatal(err)
	}

	// Load account
	accountVault, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		log.Fatal(err)
	}

	info, err := accountVault.GetAccountInfo(*fromAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PRIV:")
	fmt.Printf("%s\n", info.PrivKey)

	fmt.Println("PUB:")
	fmt.Printf("%s\n", info.PubKey)
	fmt.Println("")

	ts, err := core.GenerateJWTToken(fromAddr.Hash(), info.PrivKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ts)

	token, err := core.ValidateJWTToken(ts, fromAddr.Hash(), info.PubKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", token.Claims)
}
