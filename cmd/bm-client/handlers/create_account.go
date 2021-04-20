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
	"net"
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

type Args struct {
	Vault        *vault.Vault     // Vault
	AddrStr      string           // Textual representation of the address
	Name         string           // Name of the account
	Token        string           // Additional token for mail server
	KeyType      bmcrypto.KeyType // Type of key to generate
	RedirAddrStr string           // Textual representation of the redirect address
}

type SpinnerContext struct {
	Args            Args               // Incoming arguments
	Addr            *address.Address   // Address to create
	AccountInfo     *vault.AccountInfo // Account info found in the vault
	Proof           pow.ProofOfWork    // Generated Proof of work
	KeyPair         *bmcrypto.KeyPair   // Generated Keypair
	Reserved        bool               // True if this address is a reserved address
	ReservedDomains []string           // domains for reserved validation (if any)
	RedirAddr       *address.Address   // Address to redirect to
}

// CreateAccount creates a new account locally in the vault, stores it on the mail server and pushes the public key to the resolver
func CreateAccount(v *vault.Vault, addrStr, name, tokenStr string, kt bmcrypto.KeyType, redirStr string) {
	s := stepper.New()

	// Setup context for spinner
	spinnerCtx := SpinnerContext{
		Args: Args{
			Vault:        v,
			AddrStr:      addrStr,
			Name:         name,
			Token:        tokenStr,
			KeyType:      kt,
			RedirAddrStr: redirStr,
		},
	}
	s.Ctx = context.WithValue(s.Ctx, internal.CtxSpinnerContext, spinnerCtx)

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
		SkipIfFunc: func(s stepper.Stepper) bool { return getSpinnerContext(s).RedirAddr == nil },
		RunFunc:    checkLinkedAddressInResolver,
	})

	s.AddStep(stepper.Step{
		Title:      "Checking if token is valid and extracting data",
		SkipIfFunc: func(s stepper.Stepper) bool { return getSpinnerContext(s).Args.Token == "" },
		RunFunc:    checkToken,
	})

	s.AddStep(stepper.Step{
		Title:   "Checking if the account is already present in the vault",
		RunFunc: checkAccountInVault,
	})

	s.AddStep(stepper.Step{
		Title:          "Generating your initial keypair",
		DisplaySpinner: true,
		SkipIfFunc:     func(s stepper.Stepper) bool { return getSpinnerContext(s).AccountInfo == nil },
		RunFunc:        generateKeyPair,
	})

	s.AddStep(stepper.Step{
		Title:          fmt.Sprintf("Doing some work to let people know this is not a fake account, %sthis might take a while%s...", stepper.AnsiFgYellow, stepper.AnsiReset),
		DisplaySpinner: true,
		SkipIfFunc:     func(s stepper.Stepper) bool { return getSpinnerContext(s).AccountInfo == nil },
		RunFunc:        doProofOfWork,
	})

	s.AddStep(stepper.Step{
		Title:      "Placing your new account into the vault",
		SkipIfFunc: func(s stepper.Stepper) bool { return getSpinnerContext(s).AccountInfo == nil },
		RunFunc:    addAccountToVault,
	})

	s.AddStep(stepper.Step{
		Title:          "Checking domains for reservation proof",
		RunFunc:        checkReservedDomains,
		SkipIfFunc:     func(s stepper.Stepper) bool { return !getSpinnerContext(s).Reserved },
		DisplaySpinner: true,
	})

	s.AddStep(stepper.Step{
		Title:          "Sending your account to the server",
		DisplaySpinner: true,
		SkipIfFunc:     func(s stepper.Stepper) bool { return getSpinnerContext(s).Args.Token == "" },
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

	info := getSpinnerContext(*s).AccountInfo
	kp := info.GetActiveKey().KeyPair
	mnemonic := bminternal.WordWrap(bmcrypto.GetMnemonic(&kp), 78)

	fmt.Println(internal.GenerateFromMnemonicTemplate(internal.AccountCreatedTemplate, mnemonic))
}


func verifyAddress(s *stepper.Stepper) stepper.StepResult {
	var err error

	sc := getSpinnerContext(*s)
	sc.Addr, err = address.NewAddress(sc.Args.AddrStr)
	if err != nil {
		return s.Failure("it seems that this is not a valid address")
	}

	return s.Success("")
}

func checkReservedAddress(s *stepper.Stepper) stepper.StepResult {
	addr := getSpinnerContext(*s).Addr
	if addr == nil {
		return s.Failure("Could not find address")
	}

	ks := container.Instance.GetResolveService()
	domains, _ := ks.CheckReserved(addr.Hash())

	getSpinnerContext(*s).Reserved = len(domains) > 0
	getSpinnerContext(*s).ReservedDomains = domains

	if len(domains) > 0 {
		return s.Notice("Yes. DNS verification is needed in order to register this name")
	}

	return s.Success("Not reserved");
}

func checkReservedDomains(s *stepper.Stepper) stepper.StepResult {
	var kp *bmcrypto.KeyPair

	info := getSpinnerContext(*s).AccountInfo
	if info != nil {
		k := info.GetActiveKey().KeyPair
		kp = &k
	} else {
		kp = getSpinnerContext(*s).KeyPair
	}

	for _, domain := range getSpinnerContext(*s).ReservedDomains {
		// Check domain
		entries, err := net.LookupTXT("_bitmaelum." + domain)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fmt.Println(" --> " + entry)
			if entry == kp.PubKey.Fingerprint() {
				return s.Success("found reservation at " + domain)
			}
		}
	}

	msg := internal.GenerateFromFingerprintTemplate(internal.AccountProofTemplate, kp.PubKey.Fingerprint(), domains)
	return s.Failure(msg)
}



func checkLinkedAddressInResolver(s *stepper.Stepper) stepper.StepResult {
	redirAddr, err := address.NewAddress(getSpinnerContext(*s).Args.RedirAddrStr)
	if err != nil {
		return s.Failure("it seems that target is not a valid address")
	}

	// Store address into context
	s.Ctx = context.WithValue(s.Ctx, ctxRedir, redirAddr)


	ks := container.Instance.GetResolveService()
	_, err = ks.ResolveAddress(redirAddr.Hash())
	if err != nil {
		return s.Failure("address not found")
	}

	return s.Success("address is found")
}

func checkAddressInResolver(s *stepper.Stepper) stepper.StepResult {
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	ks := container.Instance.GetResolveService()
	_, err := ks.ResolveAddress(addr.Hash())

	if err == nil {
		return s.Failure("address already found")
	}

	return s.Success("")
}

func checkToken(s *stepper.Stepper) stepper.StepResult {
	tokenStr := s.Ctx.Value(ctxTokenStr).(string)
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	token, err := signature.ParseInviteToken(tokenStr)
	if err != nil {
		return s.Failure("it seems that this token is invalid")
	}

	// Check address matches the one in the token
	if token.AddrHash.String() != addr.Hash().String() {
		return s.Failure(fmt.Sprintf("this token is not for %s", addr.String()))
	}

	s.Ctx = context.WithValue(s.Ctx, ctxToken, token)

	return s.Success("")
}

func checkAccountInVault(s *stepper.Stepper) stepper.StepResult {
	v := s.Ctx.Value(ctxVault).(*vault.Vault)
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	if !v.HasAccount(*addr) {
		return s.Success("not found. That's good.")
	}

	info, err := v.GetAccountInfo(*addr)
	if err != nil {
		return s.Failure("found. But error while fetching from the vault.")
	}

	sc := getSpinnerContext(*s)
	sc.AccountInfo = info

	return s.Success("found. That's odd, but let's continue...")
}

func generateKeyPair(s *stepper.Stepper) stepper.StepResult {
	kt := getSpinnerContext(*s).Args.KeyType

	kp, err := bmcrypto.GenerateKeypairWithRandomSeed(kt)
	if err != nil {
		return s.Failure(err.Error())
	}

	sc := getSpinnerContext(*s)
	sc.KeyPair = kp

	return s.Success("")
}

func doProofOfWork(s *stepper.Stepper) stepper.StepResult {
	addr := s.Ctx.Value(ctxAddr).(*address.Address)

	// Find the number of bits for address creation
	res := container.Instance.GetResolveService()
	resolverCfg := res.GetConfig()

	proof := pow.NewWithoutProof(resolverCfg.ProofOfWork.Address, addr.Hash().String())
	proof.WorkMulticore()

	sc := getSpinnerContext(*s)
	sc.Proof = proof

	return s.Success("")
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
		return s.Failure(fmt.Sprintf("error while saving account into vault: %#v", err))
	}

	sc := getSpinnerContext(*s)
	sc.AccountInfo = info

	return s.Success("")
}

func uploadAccountToServer(s *stepper.Stepper) stepper.StepResult {
	token := s.Ctx.Value(ctxToken).(*signature.Token)
	info := s.Ctx.Value(ctxInfo).(*vault.AccountInfo)
	tokenStr := s.Ctx.Value(ctxTokenStr).(string)

	// Fetch routing info
	ks := container.Instance.GetResolveService()
	routingInfo, err := ks.ResolveRouting(token.RoutingID)
	if err != nil {
		return s.Failure(fmt.Sprintf("cannot find route ID inside the resolver: %#v", err))
	}

	client, err := api.NewAuthenticated(*info.Address, info.GetActiveKey().PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
	if err != nil {
		return s.Failure("error while authenticating to the API")
	}

	err = client.CreateAccount(*info, tokenStr)
	if err != nil {
		if err.Error() == "account already exists" {
			return s.Notice("account already exists on the server.")
		}

		// Other error
		return s.Failure(fmt.Sprintf("error while uploading the account: " + err.Error()))
	}

	return s.Success("")
}

func uploadAccountToResolver(s *stepper.Stepper) stepper.StepResult {
	info := getSpinnerContext(*s).AccountInfo
	if info == nil {
		return s.Failure("error while fetching account info")
	}

	ks := container.Instance.GetResolveService()
	err := ks.UploadAddressInfo(*info)
	if err != nil {
		return s.Failure(fmt.Sprintf("error while uploading account to the resolver: %s", err.Error()))
	}

	return s.Success("")
}

// getSpinnerContext returns the spinner context structure with all information that is communicated between spinner steps
func getSpinnerContext(s stepper.Stepper) *SpinnerContext {
	return s.Ctx.Value(internal.CtxSpinnerContext).(*SpinnerContext)
}
