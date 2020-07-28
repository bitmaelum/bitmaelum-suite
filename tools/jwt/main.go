package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
	Addr     string `short:"a" long:"address" description:"account"`
	Password string `short:"p" long:"password" description:"Password to decrypt your account"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	// Convert strings into addresses
	fromAddr, err := address.New(opts.Addr)
	if err != nil {
		logrus.Fatal(err)
	}

	// Load account
	accountVault, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		logrus.Fatal(err)
	}

	info, err := accountVault.GetAccountInfo(*fromAddr)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println("PRIV:")
	fmt.Printf("%s\n", info.PrivKey)

	fmt.Println("PUB:")
	fmt.Printf("%s\n", info.PubKey)
	fmt.Println("")

	privKey, err := encrypt.PEMToPrivKey([]byte(info.PrivKey))
	if err != nil {
		logrus.Fatal(err)
	}
	pubKey, err := encrypt.PEMToPubKey([]byte(info.PubKey))
	if err != nil {
		logrus.Fatal(err)
	}

	ts, err := internal.GenerateJWTToken(fromAddr.Hash(), privKey)
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(ts)
	fmt.Println("")

	token, err := internal.ValidateJWTToken(ts, fromAddr.Hash(), pubKey)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("%#v\n", token.Claims)
}
