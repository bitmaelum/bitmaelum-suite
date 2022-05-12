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
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/spf13/cobra"
)

var redirectAccountCmd = &cobra.Command{
	Use:   "redirect",
	Short: "Redirect account to another account",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		addr, err := address.NewAddress(*lnkAccount)
		if err != nil {
			fmt.Println("error: bad address: ", *lnkAccount)
			os.Exit(1)
		}

		targetAddr, err := address.NewAddress(*lnkTargetAccount)
		if err != nil {
			fmt.Println("error: bad address: ", *lnkTargetAccount)
			os.Exit(1)
		}

		info, err := v.GetAccountInfo(*addr)
		if err != nil {
			fmt.Println("error: account not found: ", addr)
			os.Exit(1)
		}

		// Check if target exists
		rs := container.Instance.GetResolveService()
		_, err = rs.ResolveAddress(targetAddr.Hash())
		if err != nil {
			fmt.Printf("error: target address '%s' not found.\n", targetAddr.String())
			os.Exit(1)
		}

		if info.RedirAddress != nil {
			fmt.Printf("Warning: you already redirect this address to '%s'. Are you sure you want to remove this and redirect to '%s' instead?\n", info.RedirAddress.String(), targetAddr.String())
		}

		fmt.Println("Linked")
	},
}

var (
	lnkAccount       *string
	lnkTargetAccount *string
)

func init() {
	accountCmd.AddCommand(redirectAccountCmd)

	lnkAccount = redirectAccountCmd.Flags().String("account", "", "Account to redirect")
	lnkTargetAccount = redirectAccountCmd.Flags().String("target", "", "Target account to redirect to")

	_ = redirectAccountCmd.MarkFlagRequired("account")
	_ = redirectAccountCmd.MarkFlagRequired("target")
}
