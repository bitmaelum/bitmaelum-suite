package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
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

	// svc := container.GetResolveService()
	// ri, err := svc.ResolveRouting(routingID)
	// if err != nil {
	// 	fmt.Println("Cannot find the specified routing ID on the resolver")
	// 	os.Exit(1)
	// }
	//
	// spew.Dump(oi)
	// spew.Dump(ri)

	hashAddr, err := address.New(addr)
	if err != nil {
		fmt.Printf("Doesn't seem like '%s' is a valid BitMealum address", addr)
		os.Exit(1)
	}
	if !hashAddr.IsOrganisationAddress() {
		fmt.Printf("Doesn't seem like '%s' is not a BitMealum organisation address", addr)
		os.Exit(1)
	}

	if hashAddr.Org != orgName {
		fmt.Printf("Address '%s' does not match your organisation '%s'", addr, orgName)
		os.Exit(1)
	}

	validUntil := time.Now().Add(7 * 24 * time.Hour)
	token, err := organisation.GenerateInviteToken(hashAddr, routingID, validUntil, oi.PrivKey)
	if err != nil {
		fmt.Println("Error while generating token: ", err)
		os.Exit(1)
	}

	fmt.Printf("Token: " + token)
}
