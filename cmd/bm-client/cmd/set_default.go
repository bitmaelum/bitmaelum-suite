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

package cmd

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var setDefaultCmd = &cobra.Command{
	Use:     "set-default",
	Aliases: []string{"default"},
	Short:   "Set default account to send from",
	Long:    `When you don't specify a from address, bm-client (ano other tools) will automatically use this address.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()

		fromAddr, err := address.NewAddress(*sdFrom)
		if err != nil {
			logrus.Fatal(err)
		}
		if !v.HasAccount(*fromAddr) {
			logrus.Fatal("Address not found in vault")
		}

		for i := range v.Store.Accounts {
			v.Store.Accounts[i].Default = false

			if v.Store.Accounts[i].Address.String() == *sdFrom {
				v.Store.Accounts[i].Default = true
			}
		}

		err = v.WriteToDisk()
		if err != nil {
			logrus.Fatalf("error while saving vault: %s\n", err)
		}

		fmt.Printf("Default address set to '%s'\n", *sdFrom)
	},
}

var sdFrom *string

func init() {
	rootCmd.AddCommand(setDefaultCmd)

	sdFrom = setDefaultCmd.Flags().String("address", "", "Default address to set")

	_ = setDefaultCmd.MarkFlagRequired("address")
}
