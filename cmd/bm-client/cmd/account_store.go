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
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/spf13/cobra"
)

var accountStoreCmd = &cobra.Command{
	Use:   "store",
	Short: "Store management",
}

var (
	astAccount *string
)

func authenticate(account string) (*api.API, *vault.AccountInfo, error) {
	v := vault.OpenDefaultVault()

	info, err := vault.GetAccount(v, *astAccount)
	if err != nil {
		return nil, nil, errors.New("cannot find account in vault")
	}

	resolver := container.Instance.GetResolveService()
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		return nil, nil, errors.New("cannot resolve routing")
	}

	client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
	if err != nil {
		return nil, nil, errors.New("cannot connect to API")
	}

	return client, info, nil
}

func init() {
	accountCmd.AddCommand(accountStoreCmd)

	astAccount = accountStoreCmd.PersistentFlags().String("account", "", "Account to set on")

	_ = accountStoreCmd.MarkPersistentFlagRequired("account")
}
