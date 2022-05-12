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

	"github.com/spf13/cobra"
)

var accountStorePutCmd = &cobra.Command{
	Use:   "put",
	Short: "Store contents into the store",
	Run: func(cmd *cobra.Command, args []string) {
		client, info, err := authenticate(*astAccount)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = client.StorePutValue(*info.StoreKey, info.Address.Hash(), *aspPath, *aspValue)
		if err != nil {
			fmt.Println("error while setting store path")
			os.Exit(1)
		}

		fmt.Println("value stored")
	},
}

var (
	aspPath  *string
	aspValue *string
)

func init() {
	accountStoreCmd.AddCommand(accountStorePutCmd)

	aspPath = accountStorePutCmd.PersistentFlags().String("path", "", "Path to store on")
	aspValue = accountStorePutCmd.PersistentFlags().String("value", "", "Value to store")

	_ = accountStorePutCmd.MarkPersistentFlagRequired("path")
	_ = accountStorePutCmd.MarkPersistentFlagRequired("value")
}
