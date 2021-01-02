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
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Create or update a specific setting in your vault",
	Long:  `Your vault accounts can have additional settings. With this command you can easily manage these.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()
		info := vault.GetAccountOrDefault(v, *sAddress)

		var msg string

		k := strings.ToLower(*sKey)
		switch k {
		case "name":
			info.Name = *sValue
			msg = fmt.Sprintf("Updated setting %s to %s", "name", *sValue)
		default:
			if *sValue == "" {
				if info.Settings != nil {
					delete(info.Settings, k)
				}
				msg = fmt.Sprintf("Removed setting %s", k)
			} else {
				if info.Settings == nil {
					info.Settings = make(map[string]string)
				}
				info.Settings[k] = *sValue
				msg = fmt.Sprintf("Updated setting %s to %s", k, *sValue)
			}
		}

		err := v.WriteToDisk()
		if err != nil {
			logrus.Fatalf("error while saving vault: %s\n", err)
		}

		logrus.Printf("%s\n", msg)
	},
}

var (
	sAddress *string
	sKey     *string
	sValue   *string
)

func init() {
	rootCmd.AddCommand(setCmd)

	sAddress = setCmd.Flags().String("address", "", "Default address to set")
	sKey = setCmd.Flags().String("key", "", "Key to set")
	sValue = setCmd.Flags().String("value", "", "Value to set")
}
