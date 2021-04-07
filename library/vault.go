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

import "github.com/bitmaelum/bitmaelum-suite/internal/vault"

// OpenVault ...
func (b *BitMaelumClient) OpenVault(path, password string) (interface{}, error) {
	var err error

	b.client.Vault, err = vault.Open(path, password)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(b.client.Vault.Store.Accounts))

	for i, acc := range b.client.Vault.Store.Accounts {
		puk := acc.GetActiveKey().PubKey
		pubkey := puk.String()

		result[i] = map[string]interface{}{
			"address":    acc.Address.String(),
			"name":       acc.Name,
			"public_key": pubkey,
		}
	}

	return result, nil
}
