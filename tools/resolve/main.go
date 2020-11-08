// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

	svc := container.Instance.GetResolveService()
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
	svc := container.Instance.GetResolveService()
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

	svc := container.Instance.GetResolveService()
	info, err := svc.ResolveOrganisation(orgHash)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("Org Name : %s\n", orgAddr)
	fmt.Printf("Hash     : %s\n", info.Hash)
	fmt.Printf("Pub key  : %s\n", info.PublicKey.String())
}
