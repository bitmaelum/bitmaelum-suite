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

package handlers

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/key"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// CreateAuthorizedKey creates a new authorized key
func CreateAuthorizedKey(info *internal.AccountInfo, targetKey *bmcrypto.PubKey, validUntil time.Duration, desc string) error {
	// Sign key
	signature, err := signKey(info.PrivKey, targetKey)
	if err != nil {
		return err
	}

	// Send key
	client, err := getAPIClient(info)
	if err != nil {
		return err
	}

	var expiry = time.Time{}
	if validUntil > 0 {
		expiry = time.Now().Add(validUntil)
	}

	key := key.NewAuthKey(info.AddressHash(), targetKey, signature, expiry, desc)
	return client.CreateAuthKey(info.AddressHash(), key)
}

func signKey(privkey bmcrypto.PrivKey, pubkey *bmcrypto.PubKey) (string, error) {
	h := hash.New(pubkey.String())

	signedKey, err := bmcrypto.Sign(privkey, h.Byte())
	if err != nil {
		return "", errors.New("cannot sing the public key")
	}

	return base64.StdEncoding.EncodeToString(signedKey), nil
}

func getAPIClient(info *internal.AccountInfo) (*api.API, error) {
	resolver := container.GetResolveService()
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		return nil, errors.New("cannot find routing ID for this account")
	}

	return api.NewAuthenticated(info, api.ClientOpts{
		Host:          routingInfo.Routing,
		AllowInsecure: config.Client.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
	})
}
