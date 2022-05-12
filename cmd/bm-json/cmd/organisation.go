// Copyright (c) 2022 BitMaelum Authors
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

package cmd

import (
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/output"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var organisationCmd = &cobra.Command{
	Use:   "organisation",
	Short: "Returns local organisation info",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := vault.Open(vault.VaultPath, vault.VaultPassword)
		if err != nil {
			output.JSONErrorStrOut("cannot open vault")
			os.Exit(1)
		}

		out := []output.JSONT{}
		for _, org := range v.Store.Organisations {

			privkey := ""
			if *orgDisplayPrivKey {
				pk := org.GetActiveKey().PrivKey
				privkey = pk.String()
			}

			pk := org.GetActiveKey().PubKey
			out = append(out, output.JSONT{
				"address":       org.Addr,
				"full_name":     org.FullName,
				"private_key":   privkey,
				"public_key":    pk.String(),
				"proof_of_work": org.Pow.String(),
				"validations":   org.Validations,
			})
		}

		output.JSONOut(out)
	},
}

var orgDisplayPrivKey *bool

func init() {
	rootCmd.AddCommand(organisationCmd)

	orgDisplayPrivKey = organisationCmd.Flags().Bool("display-private-key", false, "Should the output return the private keys as well?")
}
