package handlers

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/organisation"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	pow "github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"os"
)

// CreateOrganisation creates a new organisation locally in the vault and pushes the public key to the resolver
func CreateOrganisation(vault *vault.Vault, orgName, fullName string, orgValidations []string) {
	fmt.Printf("* Verifying if organisation name is valid: ")
	orgAddr, err := address.NewOrgHash(orgName)
	if err != nil {
		fmt.Printf("not a valid organisation")
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("ok\n")

	fmt.Printf("* Checking if your validations are correct: ")
	val, err := organisation.NewValidationTypeFromStringArray(orgValidations)
	if err != nil {
		fmt.Print("\n  X it seems that one of your validations is wrong: ", err)
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("ok.\n")

	fmt.Printf("* Checking if organisation is already known in the resolver service: ")
	ks := container.GetResolveService()
	_, err = ks.ResolveOrganisation(*orgAddr)
	if err == nil {
		fmt.Printf("\n  X it seems that this organisation is already in use. Please specify another organisation.")
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("not found. This is a good thing.\n")

	fmt.Printf("* Checking if the organisation is already present in the vault: ")
	var info *internal.OrganisationInfo
	if vault.HasOrganisation(*orgAddr) {
		fmt.Printf("\n  X organisation already present in the vault.\n")
		fmt.Println("")
		os.Exit(1)
	} else {
		fmt.Printf("not found. This is a good thing.\n")

		fmt.Printf("* Generating organisation public/private key pair: ")
		privKey, pubKey, err := bmcrypto.GenerateKeyPair(bmcrypto.KeyTypeRSA)
		if err != nil {
			fmt.Print(err)
			fmt.Println("")
			os.Exit(1)
		}
		fmt.Printf("done.\n")

		fmt.Printf("* Doing some work to let people know this is not a fake account: ")
		proof := pow.NewWithoutProof(config.Client.Accounts.ProofOfWork, orgAddr.String())
		proof.WorkMulticore()
		fmt.Printf("done.\n")

		fmt.Printf("* Adding your new organisation into the vault: ")
		info = &internal.OrganisationInfo{
			Addr:        orgName,
			Name:        fullName,
			PrivKey:     *privKey,
			PubKey:      *pubKey,
			Pow:         *proof,
			Validations: val,
		}

		vault.AddOrganisation(*info)
		err = vault.WriteToDisk()
		if err != nil {
			fmt.Printf("\n  X error while saving organisation into vault: %#v", err)
			fmt.Println("")
			os.Exit(1)
		}
		fmt.Printf("done\n")
	}

	fmt.Printf("* Making your organisation known to the outside world: ")
	err = ks.UploadOrganisationInfo(*info)
	if err != nil {
		// We can't remove the account from the vault as we have created it on the mail-server

		fmt.Printf("\n  X error while uploading organisation to the resolver: " + err.Error())
		fmt.Printf("\n  X Please try again with:\n   bm-client push-organisation -a '%s'\n", orgAddr.String())
		fmt.Println("")
		os.Exit(1)
	}
	fmt.Printf("done\n")

	fmt.Printf("\n")
	fmt.Printf("* All done")
}
