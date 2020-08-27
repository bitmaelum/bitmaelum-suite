package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-config/cmd"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"math/rand"
	"os"
	"time"
)

type options struct {
	Version bool `short:"v" long:"version" description:"Display version information"`
	AsJSON  bool `long:"json" description:"Return results as JSON"`
}

var opts options

func main() {
	rand.Seed(time.Now().UnixNano())

	internal.ParseOptions(&opts)
	if opts.Version {
		internal.WriteVersionInfo("BitMaelum Server", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	if !opts.AsJSON {
		fmt.Println(internal.GetASCIILogo())
	}

	_ = config.LoadClientConfigOrPass(".")
	_ = config.LoadServerConfigOrPass(".")

	cmd.Execute()
}
