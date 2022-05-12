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
	"fmt"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	pkginternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/signature"
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
	svc := container.Instance.GetResolveService()
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

	validUntil := pkginternal.TimeNow().Add(7 * 24 * time.Hour)
	token, err := signature.NewInviteToken(hashAddr.Hash(), routingID, validUntil, oi.GetActiveKey().PrivKey)
	if err != nil {
		fmt.Println("Error while generating token: ", err)
		os.Exit(1)
	}

	fmt.Print("Token: ", token.String())
}
