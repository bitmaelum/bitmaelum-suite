package handlers

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

func CreateOrganisationInvite(vault *vault.Vault, orgName, addr, shortRoutingID string) {
	org, err := address.NewOrgHash(orgName)
	if err != nil {
		fmt.Println("Incorrect organisation name")
		os.Exit(1)
	}

	oi, err := vault.GetOrganisationInfo(*org)
	if err != nil {
		fmt.Println("Organisation not found in the vault")
		os.Exit(1)
	}

	routingID := vault.FindShortRoutingId(shortRoutingID)
	if routingID == "" {
		routingID = shortRoutingID
	}
	fmt.Println("Found Routing ID: ", routingID)

	svc := container.GetResolveService()
	ri, err := svc.ResolveRouting(routingID)
	if err != nil {
		fmt.Println("Cannot find the specified routing ID on the resolver")
		os.Exit(1)
	}

	spew.Dump(oi)
	spew.Dump(ri)
}
