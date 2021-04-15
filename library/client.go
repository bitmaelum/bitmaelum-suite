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

package bitmaelumClient

import (
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/pkg/errors"
)

const defaultResolver = "https://resolver.bitmaelum.com"

type user struct {
	Address     *address.Address
	Name        string
	PrivateKey  *bmcrypto.PrivKey
	Vault       *vault.Vault
	RoutingInfo *resolver.RoutingInfo
}

type BitMaelumClient struct {
	user            user
	resolverService *resolver.Service
}

func NewBitMaelumClient() *BitMaelumClient {
	config.Client.Resolver.Remote.Enabled = true
	config.Client.Resolver.Remote.URL = defaultResolver

	return &BitMaelumClient{
		resolverService: container.Instance.GetResolveService(),
	}
}

func (b *BitMaelumClient) SetResolver(url string) {
	config.Client.Resolver.Remote.Enabled = true
	config.Client.Resolver.Remote.URL = url
	b.resolverService = container.Instance.GetResolveService()
}

func (b *BitMaelumClient) SetClientFromVault(accountAddress string) error {
	if b.user.Vault == nil {
		return errors.Errorf("vault not loaded")
	}

	for _, acc := range b.user.Vault.Store.Accounts {
		if acc.Address.String() == accountAddress {
			b.user.Address = acc.Address
			b.user.Name = acc.Name
			privK := acc.GetActiveKey().PrivKey
			b.user.PrivateKey = &privK
			return b.getRoutingInfo()
		}
	}

	return errors.Errorf("account %s not found on vault", accountAddress)
}

func (b *BitMaelumClient) SetClientFromMnemonic(accountAddress, name, mnemonic string) error {
	err := b.parseAccountAndName(accountAddress, name)
	if err != nil {
		return err
	}

	// Now generate a new key from the mnemonic
	kp, err := bmcrypto.GenerateKeypairFromMnemonic(strings.ToLower(mnemonic))
	if err != nil {
		return errors.Wrap(err, "parsing mnemonic")
	}

	// Verify public key belongs to account
	err = b.verifyPublicKey(kp.PubKey)
	if err != nil {
		return errors.Wrap(err, "verifying public key")
	}

	b.user.PrivateKey = &kp.PrivKey

	return b.getRoutingInfo()
}

func (b *BitMaelumClient) SetClientFromPrivateKey(accountAddress, name, privKey string) error {
	err := b.parseAccountAndName(accountAddress, name)
	if err != nil {
		return err
	}

	// Convert privKey string to bmcrypto
	b.user.PrivateKey, err = bmcrypto.PrivateKeyFromString(privKey)
	if err != nil {
		return errors.Wrap(err, "parsing private key")
	}

	pubK, err := bmcrypto.PublicKeyFromInterface(b.user.PrivateKey.Type, b.user.PrivateKey.K)
	if err != nil {
		return errors.Wrap(err, "extracting public key")
	}

	err = b.verifyPublicKey(*pubK)
	if err != nil {
		return errors.Wrap(err, "verifying public key")
	}

	return b.getRoutingInfo()
}

func (b *BitMaelumClient) parseAccountAndName(accountAddress, name string) error {
	address, err := address.NewAddress(accountAddress)
	if err != nil {
		return errors.Wrap(err, "parsing account address")
	}

	// Verify client exists
	_, err = b.resolverService.ResolveAddress(address.Hash())
	if err != nil {
		return errors.Wrap(err, "resolving client address")
	}

	b.user.Address = address
	b.user.Name = name

	return nil
}

func (b *BitMaelumClient) verifyPublicKey(pubK bmcrypto.PubKey) error {
	clientInfo, _ := b.resolverService.ResolveAddress(b.user.Address.Hash())
	if clientInfo.PublicKey.S != pubK.S {
		return errors.New("public key mismatch")
	}

	return nil
}

func (b *BitMaelumClient) getRoutingInfo() error {
	senderInfo, err := b.resolverService.ResolveAddress(b.user.Address.Hash())
	if err != nil {
		return err
	}

	b.user.RoutingInfo, err = b.resolverService.ResolveRouting(senderInfo.RoutingID)
	if err != nil {
		return err
	}

	return nil
}
