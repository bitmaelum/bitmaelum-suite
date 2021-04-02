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
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"text/template"

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
	ctxDomains
	ctxReserved
	ctxRedir
	ctxRedirStr
)

// CreateAccount creates a new account locally in the vault, stores it on the mail server and pushes the public key to the resolver
func CreateAccount(v *vault.Vault, addrStr, name, tokenStr string, kt bmcrypto.KeyType, targetStr string) {
	s := stepper.New()

	// Set some initial values in the context. We read and write to the context to deal with variables instead of using globals.
	s.Ctx = context.WithValue(s.Ctx, ctxVault, v)
	s.Ctx = context.WithValue(s.Ctx, ctxAddrStr, addrStr)
	s.Ctx = context.WithValue(s.Ctx, ctxName, name)
	s.Ctx = context.WithValue(s.Ctx, ctxTokenStr, tokenStr)
	s.Ctx = context.WithValue(s.Ctx, ctxKeyType, kt)
	s.Ctx = context.WithValue(s.Ctx, ctxRedirStr, targetStr)

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
		Title:   "Checking if address is a reserved name",
		RunFunc: checkReservedAddress,
	})

	s.AddStep(stepper.Step{
		Title:      "Checking if linked account exists",
		OnlyIfFunc: func(s stepper.Stepper) bool { return s.Ctx.Value(ctxRedirStr) != nil },
		RunFunc:    checkLinkedAddressInResolver,
	})

	s.AddStep(stepper.Step{
		Title:      "Checking if token is valid and extracting data",
		OnlyIfFunc: func(s stepper.Stepper) bool { return s.Ctx.Value(ctxTokenStr) != nil },
		RunFunc:    checkToken,
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
		Title:          "Checking domains for reservation proof",
		RunFunc:        checkReservedDomains,
		OnlyIfFunc:     func(s stepper.Stepper) bool { return s.Ctx.Value(ctxReserved) == false },
		DisplaySpinner: true,
	})

	s.AddStep(stepper.Step{
		Title:          "Sending your account to the server",
		DisplaySpinner: true,
		OnlyIfFunc:     func(s stepper.Stepper) bool { return s.Ctx.Value(ctxTokenStr) != nil },
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

	fmt.Print(`
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your account. 
If, for any reason, you lose this key, you will need to use the following 
words in order to recreate the key:
	
`)
	info := s.Ctx.Value(ctxInfo).(*vault.AccountInfo)
	kp := info.GetActiveKey().KeyPair
	fmt.Print(bminternal.WordWrap(bmcrypto.GetMnemonic(&kp), 78))
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

func checkReservedAddress(s *stepper.Stepper) stepper.StepResult {
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	ks := container.Instance.GetResolveService()
	domains, _ := ks.CheckReserved(addr.Hash())

	s.Ctx = context.WithValue(s.Ctx, ctxReserved, len(domains) > 0)
	s.Ctx = context.WithValue(s.Ctx, ctxDomains, domains)

	if len(domains) > 0 {
		return stepper.StepResult{
			Status:  stepper.NOTICE,
			Message: "Yes. DNS verification is needed in order to register this name",
		}
	}

	return stepper.StepResult{
		Status:  stepper.SUCCESS,
		Message: "Not reserved",
	}
}

func checkReservedDomains(s *stepper.Stepper) stepper.StepResult {
	var kp *bmcrypto.KeyPair

	af := s.Ctx.Value(ctxAccountFound) != nil
	if af {
		info := s.Ctx.Value(ctxInfo).(*vault.AccountInfo)
		k := info.GetActiveKey().KeyPair
		kp = &k
	} else {
		kp = s.Ctx.Value(ctxKeyPair).(*bmcrypto.KeyPair)
	}

	domains := s.Ctx.Value(ctxDomains).([]string)

	for _, domain := range domains {
		// Check domain
		entries, err := net.LookupTXT("_bitmaelum." + domain)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fmt.Println(" --> " + entry)
			if entry == kp.PubKey.Fingerprint() {
				return stepper.StepResult{
					Status:  stepper.SUCCESS,
					Message: "found reservation at " + domain,
				}
			}
		}
	}

	messageTemplate := `could not find proof in the DNS.

In order to register this reserved address, make sure you add the following information to the DNS:

    _bitmaelum TXT {{ .Fingerprint }}

This entry could be added to any of the following domains: {{ .Domains }}. Once we have found the entry, we can 
register the account onto the keyserver. For more information, please visit https://bitmaelum.com/reserved
`

	msg := generateFromTemplate(messageTemplate, kp.PubKey.Fingerprint(), domains)

	return stepper.StepResult{
		Status:  stepper.FAILURE,
		Message: msg,
	}
}

func generateFromTemplate(messageTemplate string, fingerprint string, domains []string) string {
	type tplData struct {
		Fingerprint string
		Domains     []string
	}

	data := tplData{
		Fingerprint: fingerprint,
		Domains:     domains,
	}

	msg := fmt.Sprintf("%v", data) // when things fail
	tmpl, err := template.New("template").Parse(messageTemplate)
	if err != nil {
		return msg
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return msg
	}

	return buf.String()
}

func checkLinkedAddressInResolver(s *stepper.Stepper) stepper.StepResult {
	addrStr := s.Ctx.Value(ctxRedirStr).(string)
	addr, err := address.NewAddress(addrStr)
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: "it seems that target is not a valid address",
		}
	}

	// Store address into context
	redirAddr := s.Ctx.Value(ctxRedirStr).(*address.Address)
	s.Ctx = context.WithValue(s.Ctx, ctxRedir, redirAddr)

	ks := container.Instance.GetResolveService()
	_, err = ks.ResolveAddress(addr.Hash())
	if err != nil {
		return stepper.StepResult{
			Status:  stepper.FAILURE,
			Message: "address not found",
		}
	}

	return stepper.StepResult{
		Status:  stepper.SUCCESS,
		Message: "address was found",
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
	kp, err := bmcrypto.GenerateKeypairWithRandomSeed(kt)
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

	var targetAddr *address.Address
	tmp, ok := s.Ctx.Value(ctxRedir).(*address.Address)
	if ok {
		targetAddr = tmp
	}

	var routingID string
	token, ok := s.Ctx.Value(ctxToken).(*signature.Token)
	if ok {
		routingID = token.RoutingID
	}

	info := &vault.AccountInfo{
		Address:      addr,
		RedirAddress: targetAddr,
		Name:         name,
		Keys: []vault.KeyPair{
			{
				KeyPair: *kp,
				Active:  true,
			},
		},
		Pow:       proof,
		RoutingID: routingID,
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
	info := s.Ctx.Value(ctxInfo).(*vault.AccountInfo)

	ks := container.Instance.GetResolveService()
	err := ks.UploadAddressInfo(*info)
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
