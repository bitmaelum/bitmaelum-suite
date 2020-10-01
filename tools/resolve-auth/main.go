package main

import (
	"encoding/base64"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/sirupsen/logrus"
	"os"
)

type options struct {
	Config  string `short:"c" long:"config" description:"Path to your configuration file"`
	Address string `short:"a" long:"address" description:"Address to resolve"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	v := vault.OpenVault()

	info := vault.GetAccountOrDefault(v, opts.Address)
	if info == nil {
		logrus.Fatal("No account found in vault")
		os.Exit(1)
	}

	ha, err := address.NewHash(info.Address)
	if err != nil {
		logrus.Fatal("Incorrect address")
		os.Exit(1)
	}
	hashed := []byte(ha.String() + info.RoutingID)
	sig, err := bmcrypto.Sign(info.PrivKey, hashed)
	if err != nil {
		logrus.Fatal("Cannot sign")
		os.Exit(1)
	}

	fmt.Printf("Authentication: Bearer %s\n", base64.StdEncoding.EncodeToString(sig))
}
