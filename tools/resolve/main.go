package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
)

type options struct {
	Config       string `short:"c" long:"config" description:"Path to your configuration file"`
	Address      string `short:"a" long:"address" description:"Address to resolve"`
	Organisation string `short:"o" long:"org" description:"Organisation to resolve"`
	Routing      string `short:"r" long:"routing" description:"Routing to resolve"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	if opts.Address != "" {
		resolveAddress(opts.Address)
	}
	if opts.Routing != "" {
		resolveRouting(opts.Routing)
	}
	if opts.Organisation != "" {
		resolveOrganisation(opts.Organisation)
	}
}

func resolveAddress(a string) {
	addr, err := address.NewHash(a)
	if err != nil {
		logrus.Fatal(err)
	}

	svc := container.GetResolveService()
	info, err := svc.ResolveAddress(*addr)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("BitMaelum address '%s'", a)
	fmt.Printf("Hash: %s\n", info.Hash)
	fmt.Printf("RoutingID: %s\n", info.RoutingID)
	fmt.Printf("Public key:\n%s\n", info.PublicKey.String())
}

func resolveRouting(routingID string) {
	svc := container.GetResolveService()
	info, err := svc.ResolveRouting(routingID)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("BitMaelum routing '%s'\n", routingID)
	fmt.Printf("Hash: %s\n", info.Hash)
	fmt.Printf("Public key:\n%s\n", info.PublicKey.String())
	fmt.Printf("Routing:\n%s\n", info.Routing)
}

func resolveOrganisation(org string) {
	orgHash := address.NewOrgHash(org)

	svc := container.GetResolveService()
	info, err := svc.ResolveOrganisation(orgHash)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("BitMaelum organisation '%s'\n", org)
	fmt.Printf("Hash: %s\n", info.Hash)
	fmt.Printf("Public key:\n%s\n", info.PublicKey.String())
}
