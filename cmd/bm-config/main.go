package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-config/cmd"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"os"
)

type options struct {
	Config  string `short:"c" long:"config" description:"Path to your configuration file"`
	Version bool   `short:"v" long:"version" description:"Display version information"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	if opts.Version {
		internal.WriteVersionInfo("BitMaelum Server", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println(internal.GetASCIILogo())

	config.LoadClientConfigOrPass(opts.Config)
	config.LoadServerConfigOrPass(opts.Config)

	cmd.Execute()
}
