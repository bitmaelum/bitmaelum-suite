package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/cmd"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/password"
	"os"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
	Version  bool   `short:"v" long:"version" description:"Display version information"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	if opts.Version {
		internal.WriteVersionInfo("BitMaelum Client", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println(internal.GetASCIILogo())
	config.LoadClientConfig(opts.Config)

	fromVault := false
	if opts.Password == "" {
		opts.Password, fromVault = password.AskPassword()
	}

	// Unlock vault
	accountVault, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		fmt.Printf("Error while opening vault: %s", err)
		fmt.Println("")
		os.Exit(1)
	}

	// If the password was correct and not already read from the vault, store it in the vault
	if !fromVault {
		_ = password.StorePassword(opts.Password)
	}

	// We inject our vault into the cmd, so all commands and handlers know about the vault. It doesn't
	// feel right though. I wonder if there are more idiomatic ways to do this.
	cmd.Vault = *accountVault
	cmd.Execute()
}
