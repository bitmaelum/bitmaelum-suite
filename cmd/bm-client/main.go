package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/cmd"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/vault"
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/core/password"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"os"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
	Version  bool   `short:"v" long:"version" description:"Display version information"`
}

var opts options

func main() {
	core.ParseOptions(&opts)
	if opts.Version {
		core.WriteVersionInfo("BitMaelum Client", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println(core.GetASCIILogo())
	config.LoadClientConfig(opts.Config)

	if opts.Password == "" {
		opts.Password = string(password.AskPassword())
	}

	// Unlock vault
	accountVault, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	// We inject our vault into the cmd, so all commands and handlers know about the vault. It doesn't
	// feel right though. I wonder if there are more idiomatic ways to do this.
	cmd.Vault = *accountVault
	cmd.Execute()
}
