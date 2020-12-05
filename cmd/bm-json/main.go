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
	"math/rand"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/cmd"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/output"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
	Version  bool   `short:"v" long:"version" description:"Display version information"`
}

var opts options

func main() {
	rand.Seed(time.Now().UnixNano())

	internal.ParseOptions(&opts)
	if opts.Version {
		output.JSONOut(map[string]interface{}{
			"build_date": internal.BuildDate,
			"commit_sha": internal.GitCommit,
			"version": map[string]interface{}{
				"major": internal.Version.Major,
				"minor": internal.Version.Minor,
				"patch": internal.Version.Patch,
			},
			"version_str": internal.Version.String(),
		})
		os.Exit(0)
	}

	config.LoadClientConfig(opts.Config)

	// Set main password for the vault
	vault.VaultPassword = opts.Password

	cmd.Execute()
}

