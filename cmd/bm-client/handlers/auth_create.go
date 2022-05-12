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

package handlers

import (
	"errors"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	pkginternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

var errNoRoutingID = errors.New("cannot find routing ID for this account")

// CreateAuthorizedKey creates a new authorized key
func CreateAuthorizedKey(info *vault.AccountInfo, targetKey *bmcrypto.PubKey, validUntil time.Duration, desc string) (string, error) {
	var expiry = time.Time{}
	if validUntil > 0 {
		expiry = pkginternal.TimeNow().Add(validUntil)
	}

	// Create and sign key
	k := key.NewAuthKey(info.Address.Hash(), targetKey, "", expiry, desc)
	err := k.Sign(info.GetActiveKey().PrivKey)
	if err != nil {
		return "", err
	}

	// Send key
	client, err := getAPIClient(info)
	if err != nil {
		return "", err
	}

	return client.CreateAuthKey(info.Address.Hash(), k)
}

func getAPIClient(info *vault.AccountInfo) (*api.API, error) {
	resolver := container.Instance.GetResolveService()
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		return nil, errNoRoutingID
	}

	return api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
}
