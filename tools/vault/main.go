package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/pkg/vault"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
)

type options struct {
	Config  string `short:"c" long:"config" description:"Path to your configuration file"`
	Password  string `short:"p" long:"password" description:"Password to your vault"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	v, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		panic(err)
	}

	err = v.Save()
	if err != nil {
		panic(err)
	}

	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", string(buf))
}
