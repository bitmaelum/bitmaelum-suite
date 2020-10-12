package main

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
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

func resolveAddress(addr string) {
	a, err := address.NewAddress(addr)
	if err != nil {
		logrus.Fatal(err)
	}
	addrHash := a.Hash()

	svc := container.GetResolveService()
	info, err := svc.ResolveAddress(addrHash)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("Address   : %s\n", addr)
	fmt.Printf("Hash      : %s\n", info.Hash)
	fmt.Printf("RoutingID : %s\n", info.RoutingID)
	fmt.Printf("Pub key   : %s\n", info.PublicKey.String())
}

func resolveRouting(routingID string) {
	svc := container.GetResolveService()
	info, err := svc.ResolveRouting(routingID)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("Routing ID : %s\n", routingID)
	fmt.Printf("Pub key    : %s\n", info.PublicKey.String())
	fmt.Printf("Routing    : %s\n", info.Routing)
}

func resolveOrganisation(orgAddr string) {
	orgHash := hash.New(orgAddr)

	svc := container.GetResolveService()
	info, err := svc.ResolveOrganisation(orgHash)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("Org Name : %s\n", orgAddr)
	fmt.Printf("Hash     : %s\n", info.Hash)
	fmt.Printf("Pub key  : %s\n", info.PublicKey.String())
}
