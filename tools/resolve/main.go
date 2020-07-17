package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
)

type options struct {
	Config  string `short:"c" long:"config" description:"Path to your configuration file"`
	Address string `short:"a" long:"address" description:"Address to resolve"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	addr, err := address.NewHash(opts.Address)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := container.GetResolveService()
	info, err := svc.Resolve(*addr)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("BitMaelum address '%s' resolves to %s\n", opts.Address, info.Server)
	fmt.Printf("Public key:\n%s\n", info.PublicKey)
}
