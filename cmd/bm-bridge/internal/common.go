package common

import (
	"errors"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
)

const DefaultDomain = "@bitmaelum.network"

func GetClientAndInfo(v *vault.Vault, acc string) (*vault.AccountInfo, *api.API, error) {
	//v := vault.OpenDefaultVault()

	info, err := vault.GetAccount(v, acc)
	if err != nil {
		return nil, nil, errors.New("account not found")
	}

	resolver := container.Instance.GetResolveService()
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		return nil, nil, errors.New("cannot find routing ID for this account")
	}

	client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, nil)
	if err != nil {
		return nil, nil, err
	}

	return info, client, nil
}

func AddrToEmail(address string) string {
	address = strings.Replace(address, "@", "_", -1)
	address = strings.Replace(address, "!", "", -1)

	return address + DefaultDomain
}

func EmailToAddr(email string) string {
	address := strings.Replace(email, DefaultDomain, "!", -1)
	address = strings.Replace(address, "_", "@", -1)
	return address
}
