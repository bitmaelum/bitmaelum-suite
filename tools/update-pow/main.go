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
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Password to your vault"`
	Bits     int    `short:"b" long:"bits" description:"Number of bits"`
	Force    bool   `short:"f" long:"force" description:"Force generation"`
}

func main() {
	var opts options

	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	vault.VaultPassword = opts.Password
	v := vault.OpenVault()

	for i := range v.Store.Accounts {
		var pow *proofofwork.ProofOfWork

		if len(v.Store.Accounts[i].Pow.Data) == 0 {
			var err error
			addr, err := address.NewAddress(v.Store.Accounts[i].Address)
			if err != nil {
				panic(err)
			}

			// proof of work is actually our hash address
			v.Store.Accounts[i].Pow.Data = addr.Hash().String()
		}

		if v.Store.Accounts[i].Pow.Bits >= opts.Bits && v.Store.Accounts[i].Pow.IsValid() && !opts.Force {
			fmt.Printf("Account %s has %d bits\n", v.Store.Accounts[i].Address, v.Store.Accounts[i].Pow.Bits)
			continue
		}

		fmt.Printf("Working on %s\n", v.Store.Accounts[i].Address)
		pow = proofofwork.New(opts.Bits, v.Store.Accounts[i].Pow.Data, 0)
		pow.WorkMulticore()

		v.Store.Accounts[i].Pow = pow

		err := v.WriteToDisk()
		if err != nil {
			panic(err)
		}
	}
}
