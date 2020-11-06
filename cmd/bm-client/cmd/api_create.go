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
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var apiCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create API key",
	Long:  `Your vault accounts can have additional settings. With this command you can easily manage these.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.OpenVault()
		info := vault.GetAccountOrDefault(v, *apiAddress)

		resolver := container.GetResolveService()
		routingInfo, err := resolver.ResolveRouting(info.RoutingID)
		if err != nil {
			logrus.Fatal("Cannot find routing ID for this account")
			os.Exit(1)
		}

		var expiry = time.Time{}
		fmt.Printf("EXPIRY: %#v\n", expiry)
		if *acValidUntil > 0 {
			expiry = time.Now().Add(*acValidUntil)
			fmt.Printf("UNTIL EXPIRY: %#v\n", expiry)
		}

		key := key.NewAPIKey(*acPerms, expiry, *acDesc)
		fmt.Printf("%#v\n", key)
		fmt.Printf("%#v\n", *acValidUntil)

		client, err := api.NewAuthenticated(info.Address, &info.PrivKey, routingInfo.Routing)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}

		err = client.CreateAPIKey(info.AddressHash(), key)
		if err != nil {
			logrus.Fatal(err)
			os.Exit(1)
		}
	},
}

var (
	acPerms      *[]string
	acValidUntil *time.Duration
	acDesc       *string
)

func init() {
	apiCmd.AddCommand(apiCreateCmd)

	acPerms = apiCreateCmd.Flags().StringArray("perm", []string{}, "Permissions to set")
	acValidUntil = apiCreateCmd.Flags().Duration("duration", time.Duration(0), "Time valid")
	acDesc = apiCreateCmd.Flags().StringP("desc", "d", "", "Value to set")

	_ = apiCreateCmd.MarkFlagRequired("perm")
	_ = apiCreateCmd.MarkFlagRequired("desc")
}
