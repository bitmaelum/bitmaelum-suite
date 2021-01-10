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

package vault

import (
	"encoding/json"
	"errors"

	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

// MigrateVault will migrate a vault from a specific version all the way to the latest version
func MigrateVault(data []byte, fromVersion int) (*StoreType, error) {
	storeV0 := make([]AccountInfoV1, 0)
	storeV1 := &StoreTypeV1{}
	storeV2 := &StoreType{}

	// Read the correct initial data format
	switch fromVersion {
	case VersionV0:
		err := json.Unmarshal(data, &storeV0)
		if err != nil {
			return nil, err
		}
	case VersionV1:
		err := json.Unmarshal(data, &storeV1)
		if err != nil {
			return nil, err
		}
	case VersionV2:
		err := json.Unmarshal(data, &storeV2)
		if err != nil {
			return nil, err
		}
	}

	// Iterate migrations until we reach the latest version
	for fromVersion <= LatestVaultVersion {
		switch fromVersion {
		case VersionV0:
			/*
			 * V0 -> V1 adds organisations
			 */
			storeV1.Accounts = storeV0
			storeV1.Organisations = []OrganisationInfoV1{}

			fromVersion++
		case VersionV1:
			/*
			 * V1 -> V2 has multiple keys in both organisation and accounts
			 */
			for _, accV1 := range storeV1.Accounts {
				accV2 := &AccountInfo{}
				accV2.Address = accV1.Address
				accV2.Pow = accV1.Pow
				accV2.RoutingID = accV1.RoutingID
				accV2.Name = accV1.Name
				accV2.Settings = accV1.Settings
				accV2.Keys = []KeyPair{
					{
						KeyPair: bmcrypto.KeyPair{
							Generator:   "",
							FingerPrint: accV1.PubKey.Fingerprint(),
							PrivKey:     accV1.PrivKey,
							PubKey:      accV1.PubKey,
						},
						Active: true,
					},
				}
				storeV2.Accounts = append(storeV2.Accounts, *accV2)
			}

			for _, orgV1 := range storeV1.Organisations {
				orgV2 := &OrganisationInfo{}
				orgV2.Addr = orgV1.Addr
				orgV2.Pow = orgV1.Pow
				orgV2.FullName = orgV1.FullName
				orgV2.Validations = orgV1.Validations
				orgV2.Keys = []KeyPair{
					{
						KeyPair: bmcrypto.KeyPair{
							Generator:   "",
							FingerPrint: orgV1.PubKey.Fingerprint(),
							PrivKey:     orgV1.PrivKey,
							PubKey:      orgV1.PubKey,
						},
						Active: true,
					},
				}
				storeV2.Organisations = append(storeV2.Organisations, *orgV2)
			}

			fromVersion++

		case VersionV2:
			// Latest version, no need to migrate
			return storeV2, nil
		}
	}

	return nil, errors.New("error while running migrations on the vault")
}
