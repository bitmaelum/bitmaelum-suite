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

package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal/stepper"
	bminternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/signature"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

type ctxKey int

const (
	ctxAddr ctxKey = iota
	ctxAddrStr
	ctxVault
	ctxName
	ctxTokenStr
	ctxKeyType
	ctxInfo
	ctxAccountFound
	ctxProof
	ctxToken
	ctxKeyPair
)

// CreateAccount creates a new account locally in the vault, stores it on the mail server and pushes the public key to the resolver
func CreateAccount(v *vault.Vault, addrStr, name, tokenStr string, kt bmcrypto.KeyType) {
	s := stepper.New()

	// Set some initial values in the context. We read and write to the context to deal with variables instead of using globals.
	s.Ctx = context.WithValue(s.Ctx, ctxVault, v)
	s.Ctx = context.WithValue(s.Ctx, ctxAddrStr, addrStr)
	s.Ctx = context.WithValue(s.Ctx, ctxName, name)
	s.Ctx = context.WithValue(s.Ctx, ctxTokenStr, tokenStr)
	s.Ctx = context.WithValue(s.Ctx, ctxKeyType, kt)

	// Add all the steps from the account creation procedure

	s.AddStep(stepper.Step{
		Title:   fmt.Sprintf("Verifying if address %s%s%s is correct", stepper.AnsiFgYellow, addrStr, stepper.AnsiReset),
		RunFunc: verifyAddress,
	})

	s.AddStep(stepper.Step{
		Title:   "Checking if address is already known in the resolver service",
		RunFunc: checkAddressInResolver,
	})

	s.AddStep(stepper.Step{
		Title:   "Checking if token is valid and extracting data",
		RunFunc: checkToken,
	})

	s.AddStep(stepper.Step{
		Title:   "Checking if the account is already present in the vault",
		RunFunc: checkAccountInVault,
	})

	s.AddStep(stepper.Step{
		Title:          "Generating your initial keypair",
		DisplaySpinner: true,
		OnlyIfFunc:     accountNotFoundInContext,
		RunFunc:        generateKeyPair,
	})

	s.AddStep(stepper.Step{
		Title:          fmt.Sprintf("Doing some work to let people know this is not a fake account, %sthis might take a while%s...", stepper.AnsiFgYellow, stepper.AnsiReset),
		DisplaySpinner: true,
		OnlyIfFunc:     accountNotFoundInContext,
		RunFunc:        doProofOfWork,
	})

	s.AddStep(stepper.Step{
		Title:      "Placing your new account into the vault",
		OnlyIfFunc: accountNotFoundInContext,
		RunFunc:    addAccountToVault,
	})

	s.AddStep(stepper.Step{
		Title:          "Sending your account to the server",
		DisplaySpinner: true,
		RunFunc:        uploadAccountToServer,
	})

	s.AddStep(stepper.Step{
		Title:          "Making your account known to the outside world",
		DisplaySpinner: true,
		RunFunc:        uploadAccountToResolver,
	})

	// Run the stepper
	s.Run()
	if s.Status == stepper.FAILURE {
		fmt.Println("There was an error while creating the account.")
		os.Exit(1)
	}

	info := s.Ctx.Value(ctxInfo).(*vault.AccountInfo)

	fmt.Print(`
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your account. 
If, for any reason, you lose this key, you will need to use the following 
words in order to recreate the key:
	
`)
	kp := info.GetActiveKey().KeyPair
	fmt.Print(bminternal.WordWrap(bminternal.GetMnemonic(&kp), 78))
	fmt.Print(`
	
Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.
	
WITHOUT THESE WORDS, ALL ACCESS TO YOUR ACCOUNT IS LOST!
	`)
}

func verifyAddress(s *stepper.Stepper) stepper.StepResult {
	addrStr := s.Ctx.Value(ctxAddrStr).(string)

	addr, err := address.NewAddress(addrStr)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: "it seems that this is not a valid address",
		}
	}

	s.Ctx = context.WithValue(s.Ctx, ctxAddr, addr)

	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func checkAddressInResolver(s *stepper.Stepper) stepper.StepResult {
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	ks := container.Instance.GetResolveService()
	_, err := ks.ResolveAddress(addr.Hash())

	if err == nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: "address already found",
		}
	}

	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func checkToken(s *stepper.Stepper) stepper.StepResult {
	tokenStr := s.Ctx.Value(ctxTokenStr).(string)
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	token, err := signature.ParseInviteToken(tokenStr)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: "it seems that this token is invalid",
		}
	}

	// Check address matches the one in the token
	if token.AddrHash.String() != addr.Hash().String() {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: fmt.Sprintf("this token is not for %s", addr.String()),
		}
	}

	s.Ctx = context.WithValue(s.Ctx, ctxToken, token)

	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func checkAccountInVault(s *stepper.Stepper) stepper.StepResult {
	v := s.Ctx.Value(ctxVault).(*vault.Vault)
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	if !v.HasAccount(*addr) {
		return stepper.StepResult{
			Status:  stepper.SUCCESS,
			Message: "not found. That's good.",
		}
	}

	info, err := v.GetAccountInfo(*addr)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: "found. But error while fetching from the vault.",
		}
	}

	s.Ctx = context.WithValue(s.Ctx, ctxInfo, info)
	s.Ctx = context.WithValue(s.Ctx, ctxAccountFound, true)

	return stepper.StepResult{
		Status:  stepper.SUCCESS,
		Message: "found. That's odd, but let's continue...",
	}

}

func generateKeyPair(s *stepper.Stepper) stepper.StepResult {
	kt := s.Ctx.Value(ctxKeyType).(bmcrypto.KeyType)
	kp, err := bminternal.GenerateKeypairWithRandomSeed(kt)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: err.Error(),
		}
	}

	s.Ctx = context.WithValue(s.Ctx, ctxKeyPair, kp)
	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func doProofOfWork(s *stepper.Stepper) stepper.StepResult {
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	// Find the number of bits for address creation
	res := container.Instance.GetResolveService()
	resolverCfg := res.GetConfig()

	proof := pow.NewWithoutProof(resolverCfg.ProofOfWork.Address, addr.Hash().String())
	proof.WorkMulticore()

	s.Ctx = context.WithValue(s.Ctx, ctxProof, proof)

	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func addAccountToVault(s *stepper.Stepper) stepper.StepResult {
	v := s.Ctx.Value(ctxVault).(*vault.Vault)
	addr := s.Ctx.Value(ctxAddr).(*address.Address)
	name := s.Ctx.Value(ctxName).(string)
	kp := s.Ctx.Value(ctxKeyPair).(*bmcrypto.KeyPair)
	proof := s.Ctx.Value(ctxProof).(*pow.ProofOfWork)
	token := s.Ctx.Value(ctxToken).(*signature.Token)

	info := &vault.AccountInfo{
		Address: addr,
		Name:    name,
		Keys: []vault.KeyPair{
			{
				KeyPair: *kp,
				Active:  true,
			},
		},
		Pow:       proof,
		RoutingID: token.RoutingID,
	}

	v.AddAccount(*info)

	err := v.Persist()
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: fmt.Sprintf("error while saving account into vault: %#v", err),
		}
	}

	s.Ctx = context.WithValue(s.Ctx, ctxInfo, info)
	s.Ctx = context.WithValue(s.Ctx, ctxAccountFound, true)

	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func uploadAccountToServer(s *stepper.Stepper) stepper.StepResult {
	token := s.Ctx.Value(ctxToken).(*signature.Token)
	info := s.Ctx.Value(ctxInfo).(*vault.AccountInfo)
	tokenStr := s.Ctx.Value(ctxTokenStr).(string)

	// Fetch routing info
	ks := container.Instance.GetResolveService()
	routingInfo, err := ks.ResolveRouting(token.RoutingID)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: fmt.Sprintf("cannot find route ID inside the resolver: %#v", err),
		}
	}

	client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: "error while authenticating to the API",
		}
	}

	err = client.CreateAccount(*info, tokenStr)
	if err != nil {
		if err.Error() == "account already exists" {
			return stepper.StepResult{
				Status:  stepper.NOTICE,
				Message: "account already exists on the server.",
			}
		}

		// Other error
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: fmt.Sprintf("error while uploading the account: " + err.Error()),
		}
	}

	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func uploadAccountToResolver(s *stepper.Stepper) stepper.StepResult {
	addr := s.Ctx.Value(ctxAddr).(*address.Address)
	info := s.Ctx.Value(ctxInfo).(*vault.AccountInfo)
	tokenStr := s.Ctx.Value(ctxTokenStr).(string)

	if !addr.HasOrganisationPart() {
		tokenStr = ""
	}

	ks := container.Instance.GetResolveService()
	err := ks.UploadAddressInfo(*info, tokenStr)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: fmt.Sprintf("error while uploading account to the resolver: %s", err.Error()),
		}
	}

	return stepper.StepResult{
		Status: stepper.SUCCESS,
	}
}

func accountNotFoundInContext(s stepper.Stepper) bool {
	return s.Ctx.Value(ctxAccountFound) != nil
}
