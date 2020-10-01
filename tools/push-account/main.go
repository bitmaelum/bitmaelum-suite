package main

import (
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
	Address  string `short:"a" long:"address" description:"Address" default:""`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	logrus.SetLevel(logrus.TraceLevel)

	vault.VaultPassword = opts.Password
	v := vault.OpenVault()

	info := vault.GetAccountOrDefault(v, opts.Address)
	if info == nil {
		logrus.Fatal("No account found in vault")
		os.Exit(1)
	}

	rs := container.GetResolveService()
	err := rs.UploadAddressInfo(*info)
	if err != nil {
		fmt.Printf("Error for account %s: %s\n", info.Address, err)
	}
}
