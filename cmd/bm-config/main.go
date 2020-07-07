package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-config/cmd"
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"os"
)

type options struct {
	Config  string `short:"c" long:"config" description:"Path to your configuration file"`
	Version bool   `short:"v" long:"version" description:"Display version information"`
}

var opts options

func main() {
	core.ParseOptions(&opts)
	if opts.Version {
		core.WriteVersionInfo("BitMaelum Server", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	fmt.Println(core.GetASCIILogo())

	config.LoadClientConfigOrPass(opts.Config)
	config.LoadServerConfigOrPass(opts.Config)

	cmd.Execute()
}
