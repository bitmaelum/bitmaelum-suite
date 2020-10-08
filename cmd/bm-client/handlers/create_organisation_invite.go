package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/invite"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
)

// CreateOrganisationInvite invites a user to an organisation on the given routing-id
func CreateOrganisationInvite(vault *vault.Vault, orgAddr, inviteAddr, shortRoutingID string) {
	oi, err := vault.GetOrganisationInfo(hash.New(orgAddr))
	if err != nil {
		fmt.Println("Organisation not found in the vault")
		os.Exit(1)
	}

	routingID := vault.FindShortRoutingID(shortRoutingID)
	if routingID == "" {
		routingID = shortRoutingID
	}
	fmt.Println("Found Routing ID: ", routingID)

	// Verify the routing ID exists
	svc := container.GetResolveService()
	_, err = svc.ResolveRouting(routingID)
	if err != nil {
		fmt.Println("Cannot find the specified routing ID on the resolver")
		os.Exit(1)
	}

	hashAddr, err := address.NewAddress(inviteAddr)
	if err != nil {
		fmt.Printf("Doesn't seem like '%s' is a valid BitMealum address", inviteAddr)
		os.Exit(1)
	}
	if !hashAddr.HasOrganisationPart() {
		fmt.Printf("Doesn't seem like '%s' is not a BitMealum organisation address", inviteAddr)
		os.Exit(1)
	}

	if hashAddr.Org != orgAddr {
		fmt.Printf("Address '%s' does not match your organisation '%s'", inviteAddr, orgAddr)
		os.Exit(1)
	}

	validUntil := time.Now().Add(7 * 24 * time.Hour)
	token, err := invite.NewInviteToken(hashAddr.Hash(), routingID, validUntil, oi.PrivKey)
	if err != nil {
		fmt.Println("Error while generating token: ", err)
		os.Exit(1)
	}

	fmt.Print("Token: ", token.String())
}
