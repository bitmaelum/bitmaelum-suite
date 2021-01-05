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
	"math/rand"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/cmd"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	bminternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	bminternal.ParseOptions(&internal.Opts)
	if internal.Opts.Version {
		bminternal.WriteVersionInfo("BitMaelum Client", os.Stdout)
		fmt.Println()
		os.Exit(1)
	}

	// Set default vault info if set in config
	vault.VaultPassword = internal.Opts.Password
	vault.VaultPath = config.Client.Vault.Path
	if internal.Opts.Vault != "" {
		vault.VaultPath = internal.Opts.Vault
	}

	cmd.Execute()
}
