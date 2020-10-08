package main

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
)

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Password to your vault"`
	Bits     int    `short:"b" long:"bits" description:"Number of bits"`
	Force    bool   `short:"f" long:"force" description:"Force generation"`
}

func main() {
	var opts options

	internal.ParseOptions(&opts)
	config.LoadClientConfig(opts.Config)

	vault.VaultPassword = opts.Password
	v := vault.OpenVault()

	for i := range v.Store.Accounts {
		var pow *proofofwork.ProofOfWork

		if len(v.Store.Accounts[i].Pow.Data) == 0 {
			var err error
			addr, err := address.NewAddress(v.Store.Accounts[i].Address)
			if err != nil {
				panic(err)
			}

			// proof of work is actually our hash address
			v.Store.Accounts[i].Pow.Data = addr.Hash().String()
		}

		if v.Store.Accounts[i].Pow.Bits >= opts.Bits && v.Store.Accounts[i].Pow.IsValid() && !opts.Force {
			fmt.Printf("Account %s has %d bits\n", v.Store.Accounts[i].Address, v.Store.Accounts[i].Pow.Bits)
			continue
		}

		fmt.Printf("Working on %s\n", v.Store.Accounts[i].Address)
		pow = proofofwork.New(opts.Bits, v.Store.Accounts[i].Pow.Data, 0)
		pow.WorkMulticore()

		v.Store.Accounts[i].Pow = *pow

		err := v.WriteToDisk()
		if err != nil {
			panic(err)
		}
	}
}
