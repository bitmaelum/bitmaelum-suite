package bitmaelumClient

import "github.com/bitmaelum/bitmaelum-suite/internal/vault"

// OpenVault ...
func (b *BitMaelumClient) OpenVault(path, password string) (interface{}, error) {
	v, err := vault.Open(path, password)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(v.Store.Accounts))

	for i, acc := range v.Store.Accounts {
		pk := acc.GetActiveKey().PrivKey
		privkey := pk.String()

		puk := acc.GetActiveKey().PubKey
		pubkey := puk.String()

		result[i] = map[string]interface{}{
			"address":     acc.Address.String(),
			"hash":        acc.Address.Hash().String(),
			"name":        acc.Name,
			"routing_id":  acc.RoutingID,
			"private_key": privkey,
			"public_key":  pubkey,
		}
	}

	return result, nil
}
