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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var accountSettingsSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Create or update a specific setting in your account",
	Long:  `Your accounts can have additional settings. With this command you can easily manage these.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenDefaultVault()

		info, err := vault.GetAccount(v, *asAccount)
		if err != nil {
			fmt.Println("cannot find account in vault")
			os.Exit(1)
		}

		var msg string

		k := strings.ToLower(*assKey)
		switch k {
		case "name":
			info.Name = *assValue
			msg = fmt.Sprintf("Updated setting %s to %s", "name", *assValue)
		default:
			if info.Settings == nil {
				info.Settings = make(map[string]string)
			}
			info.Settings[k] = *assValue
			msg = fmt.Sprintf("Updated setting %s to %s", k, *assValue)
		}

		err = v.WriteToDisk()
		if err != nil {
			fmt.Printf("error while saving vault: %s\n", err)
			os.Exit(1)
		}

		logrus.Printf("%s\n", msg)
	},
}

var (
	assKey     *string
	assValue   *string
)

func init() {
	accountSettingsCmd.AddCommand(accountSettingsSetCmd)

	assKey = accountSettingsSetCmd.Flags().String("key", "", "Key to set")
	assValue = accountSettingsSetCmd.Flags().String("value", "", "Value to set")

	_ = accountSettingsSetCmd.MarkFlagRequired("key")
	_ = accountSettingsSetCmd.MarkFlagRequired("value")
}
