// Copyright (c) 2021 BitMaelum Authors
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
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-imap/internal"
	bminternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
	Vault    string `long:"vault" description:"Custom vault file" default:""`
	Version  bool   `short:"v" long:"version" description:"Display version information"`
}

// Opts are the options set through the command line
var Opts options


var opts options

func main() {
	rand.Seed(time.Now().UnixNano())

	bminternal.ParseOptions(&opts)
	if opts.Version {
		bminternal.WriteVersionInfo("BitMaelum IMAP", os.Stdout)
		fmt.Println()
		os.Exit(0)
	}

	config.LoadClientConfig(opts.Config)

	// Set default vault info if set in config
	vault.VaultPassword = opts.Password
	vault.VaultPath = config.Client.Vault.Path
	if opts.Vault != "" {
		vault.VaultPath = opts.Vault
	}

	v, err := vault.Open(vault.VaultPath, vault.VaultPassword)
	if err != nil {
		fmt.Println("error: cannot open vault")
		os.Exit(1)
	}

	fmt.Println("BM-IMAP started...")

	ln, err := net.Listen("tcp", "127.0.0.1:1143")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("error while accepting connection: %s\n", err)
		}

		c := internal.NewConn(conn, v)
		go c.Handle()
	}
}
