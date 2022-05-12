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

package internal

import (
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-json/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
)

// GetClientAndInfo is a simple wrapper that will fetch the vault, account info from the given vault and an authenticated
// api client. Since this is used a lot, we created a separate function for this. This will take care of a lot of code
// duplication.
func GetClientAndInfo(acc string) (*vault.Vault, *vault.AccountInfo, *api.API, error) {
	v := vault.OpenDefaultVault()

	info, err := vault.GetAccount(v, acc)
	if err != nil {
		return nil, nil, nil, errors.New("account not found")
	}

	resolver := container.Instance.GetResolveService()
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		return nil, nil, nil, errors.New("cannot find routing ID for this account")
	}

	client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, JwtJSONErrorFunc)
	if err != nil {
		return nil, nil, nil, err
	}

	return v, info, client, nil
}
