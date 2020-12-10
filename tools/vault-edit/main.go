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
)

type options struct {
	Config      string `short:"c" long:"config" description:"Path to your configuration file"`
	Password    string `short:"p" long:"password" description:"Password to your vault"`
	NewPassword string `short:"n" long:"new-password" description:"New password to your vault"`
}

var opts options

func main() {
	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	v, err := vault.New(config.Client.Accounts.Path, []byte(opts.Password))
	if err != nil {
		panic(err)
	}

	st := vault.StoreType{}
	err = internal.JSONFileEditor(v.Store, &st)
	if err != nil {
		panic(err)
	}

	v.Store = st

	if opts.NewPassword != "" {
		v.ChangePassword(opts.NewPassword)
	}

	err = v.WriteToDisk()
	if err != nil {
		panic(err)
	}

	fmt.Println("Vault saved to disk")
}
