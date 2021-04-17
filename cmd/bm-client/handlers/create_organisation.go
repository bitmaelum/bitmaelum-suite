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
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

type OrgArgs struct {
	Vault       *vault.Vault     // Vault
	AddrStr     string           // Textual representation of the address
	Name        string           // Name of the account
	TokenStr    string           // Additional token for mail server
	KeyType     bmcrypto.KeyType // Type of key to generate
	Validations []string         // Validations
}

type OrgSpinnerContext struct {
	Args            OrgArgs            // Incoming arguments
	//Addr            *address.Address   // Address to create
	//AccountInfo     *vault.AccountInfo // Account info found in the vault
	//Proof           *pow.ProofOfWork   // Generated Proof of work
	//KeyPair         *bmcrypto.KeyPair  // Generated Keypair
	Reserved        bool               // True if this address is a reserved address
	ReservedDomains []string           // domains for reserved validation (if any)
	//Token           *signature.Token   // Information from given token
	Organisation     *vault.OrganisationInfo				// Organisation info
}

// getOrganisationContext returns the spinner context structure with all information that is communicated between spinner steps
func getOrganisationContext(s stepper.Stepper) *OrgSpinnerContext {
	return s.Ctx.Value(internal.CtxSpinnerContext).(*OrgSpinnerContext)
}

// CreateOrganisation creates a new organisation locally in the vault and pushes the public key to the resolver
func CreateOrganisation(v *vault.Vault, orgAddr, fullName string, orgValidations []string, kt bmcrypto.KeyType) {
	s := stepper.New()

	// Setup context for spinner
	spinnerCtx := OrgSpinnerContext{
		Args: OrgArgs{
			Vault:       v,
			KeyType:     kt,
			Name:        fullName,
			Validations: orgValidations,
			AddrStr:     orgAddr,
		},
	}
	s.Ctx = context.WithValue(s.Ctx, internal.CtxSpinnerContext, spinnerCtx)

	// Add all the steps from the account creation procedure

	s.AddStep(stepper.Step{
		Title:   "Checking if organisation is already known in the resolver service",
		RunFunc: checkOrganisationInResolver,
	})

	s.AddStep(stepper.Step{
		Title:   "Checking if address is a reserved name",
		RunFunc: checkOrganisationReservedAddress,
	})

	s.AddStep(stepper.Step{
		Title:          "Checking if validations are correct",
		DisplaySpinner: true,
		RunFunc:        checkValidations,
	})

	s.AddStep(stepper.Step{
		Title:   "Checking if the organisation is already present in the vault",
		RunFunc: checkOrganisationInVault,
	})

	s.AddStep(stepper.Step{
		Title:          "Generating organisation public/private keypair",
		DisplaySpinner: true,
		SkipIfFunc: func(s stepper.Stepper) bool { return getOrganisationContext(*s).Organisation == nil },
		RunFunc:        generateOrganisationKeyPair,
	})

	s.AddStep(stepper.Step{
		Title:          fmt.Sprintf("Doing some work to let people know this is not a fake account, %sthis might take a while%s...", stepper.AnsiFgYellow, stepper.AnsiReset),
		DisplaySpinner: true,
		SkipIfFunc: func(s stepper.Stepper) bool { return getOrganisationContext(*s).Organisation == nil },
		RunFunc:        doProofOfWorkOrg,
	})

	s.AddStep(stepper.Step{
		Title:      "Placing your new organisation into the vault",
		SkipIfFunc: func(s stepper.Stepper) bool { return getOrganisationContext(*s).Organisation == nil },
		RunFunc:    addOrganisationToVault,
	})

	s.AddStep(stepper.Step{
		Title:          "Checking domains for reservation proof",
		RunFunc:        checkOrganisationReservedDomains,
		SkipIfFunc:     func(s stepper.Stepper) bool { return s.Ctx.Value(ctxOrgReserved).(bool) == false },
		DisplaySpinner: true,
	})

	s.AddStep(stepper.Step{
		Title:          "Making your organisation known to the outside world",
		DisplaySpinner: true,
		RunFunc:        uploadOrganisationToResolver,
	})

	// Run the stepper
	s.Run()
	if s.Status == stepper.FAILURE {
		fmt.Println("There was an error while creating the organisation.")
		os.Exit(1)
	}

	kp := getOrganisationContext(*s).Organisation.GetActiveKey().KeyPair
	mnemonic := bminternal.WordWrap(bmcrypto.GetMnemonic(&kp), 78)

	fmt.Println(internal.GenerateFromMnemonicTemplate(internal.OrganisationCreatedTemplate, mnemonic))
}

func checkOrganisationInVault(s *stepper.Stepper) stepper.StepResult {
	v := s.Ctx.Value(ctxOrgVault).(*vault.Vault)
	orgHash := s.Ctx.Value(ctxOrgHash).(hash.Hash)

	if !v.HasOrganisation(orgHash) {
		return s.Success("not found. That's good.")
	}

	info, err := v.GetOrganisationInfo(orgHash)
	if err != nil {
		return s.Failure("found. But error while fetching from the vault.")
	}

	s.Ctx = context.WithValue(s.Ctx, ctxOrgInfo, info)
	s.Ctx = context.WithValue(s.Ctx, ctxOrganisationFound, true)

	return s.Success("found. That's odd, but let's continue...")
}

func checkOrganisationInResolver(s *stepper.Stepper) stepper.StepResult {
	orgHash := s.Ctx.Value(ctxOrgHash).(hash.Hash)

	ks := container.Instance.GetResolveService()
	_, err := ks.ResolveOrganisation(orgHash)

	if err == nil {
		return s.Failure("organisation already found")
	}

	return s.Success("")
}

func checkValidations(s *stepper.Stepper) stepper.StepResult {
	arr := s.Ctx.Value(ctxOrgValidationsStr).([]string)
	validations, err := organisation.NewValidationTypeFromStringArray(arr)
	if err != nil {
		return s.Failure("validation failed")
	}

	s.Ctx = context.WithValue(s.Ctx, ctxOrgValidations, validations)

	return s.Success("")
}

func doProofOfWorkOrg(s *stepper.Stepper) stepper.StepResult {
	orgHash := s.Ctx.Value(ctxOrgHash).(hash.Hash)

	// Find the number of bits for address creation
	res := container.Instance.GetResolveService()
	resolverCfg := res.GetConfig()

	proof := pow.NewWithoutProof(resolverCfg.ProofOfWork.Organisation, orgHash.String())
	proof.WorkMulticore()

	s.Ctx = context.WithValue(s.Ctx, ctxOrgProof, proof)

	return s.Success("")
}

func generateOrganisationKeyPair(s *stepper.Stepper) stepper.StepResult {
	kt := s.Ctx.Value(ctxOrgKeyType).(bmcrypto.KeyType)
	kp, err := bmcrypto.GenerateKeypairWithRandomSeed(kt)
	if err != nil {
		return s.Failure(err.Error())
	}

	s.Ctx = context.WithValue(s.Ctx, ctxOrgKeyPair, kp)
	return s.Success("")
}

func addOrganisationToVault(s *stepper.Stepper) stepper.StepResult {
	v := s.Ctx.Value(ctxOrgVault).(*vault.Vault)
	orgAddr := s.Ctx.Value(ctxOrgAddr).(string)
	name := s.Ctx.Value(ctxOrgName).(string)
	kp := s.Ctx.Value(ctxOrgKeyPair).(*bmcrypto.KeyPair)
	proof := s.Ctx.Value(ctxOrgProof).(*pow.ProofOfWork)
	validations := s.Ctx.Value(ctxOrgValidations).([]organisation.ValidationType)

	info := &vault.OrganisationInfo{
		Addr:     orgAddr,
		FullName: name,
		Keys: []vault.KeyPair{
			{
				KeyPair: *kp,
				Active:  true,
			},
		},
		Pow:         proof,
		Validations: validations,
	}

	v.AddOrganisation(*info)
	err := v.Persist()
	if err != nil {
		return s.Failure(fmt.Sprintf("error while saving organisation into vault: %#v", err))
	}

	s.Ctx = context.WithValue(s.Ctx, ctxOrgInfo, info)
	s.Ctx = context.WithValue(s.Ctx, ctxOrganisationFound, true)

	return s.Success("")
}

func uploadOrganisationToResolver(s *stepper.Stepper) stepper.StepResult {
	info := s.Ctx.Value(ctxOrgInfo).(*vault.OrganisationInfo)

	ks := container.Instance.GetResolveService()
	err := ks.UploadOrganisationInfo(*info)
	if err != nil {
		return s.Failure(fmt.Sprintf("error while uploading organisation to the resolver: %s", err.Error()))
	}

	return s.Success("")
}

func checkOrganisationReservedAddress(s *stepper.Stepper) stepper.StepResult {
	orgAddr := s.Ctx.Value(ctxOrgAddr).(string)
	orgHash := hash.New(orgAddr)

	ks := container.Instance.GetResolveService()
	domains, _ := ks.CheckReserved(orgHash)

	s.Ctx = context.WithValue(s.Ctx, ctxOrgReserved, len(domains) > 0)
	s.Ctx = context.WithValue(s.Ctx, ctxOrgDomains, domains)

	if len(domains) > 0 {
		return s.Notice("Yes. DNS verification is needed in order to register this organisation")
	}

	return s.Success("not reserved")
}

func checkOrganisationReservedDomains(s *stepper.Stepper) stepper.StepResult {
	var kp *bmcrypto.KeyPair

	af := s.Ctx.Value(ctxOrganisationFound) != nil
	if af {
		info := s.Ctx.Value(ctxOrgInfo).(*vault.OrganisationInfo)
		k := info.GetActiveKey().KeyPair
		kp = &k
	} else {
		kp = s.Ctx.Value(ctxOrgKeyPair).(*bmcrypto.KeyPair)
	}

	domains := s.Ctx.Value(ctxOrgDomains).([]string)

	for _, domain := range domains {
		// Check domain
		entries, err := net.LookupTXT("_bitmaelum." + domain)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry == kp.PubKey.Fingerprint() {
				return s.Success("found reservation at " + domain)
			}
		}
	}

	msg := internal.GenerateFromFingerprintTemplate(internal.OrganisationProofTemplate, kp.PubKey.Fingerprint(), domains)
	return s.Failure(msg)
}
