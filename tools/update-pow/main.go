package main

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	tools "github.com/bitmaelum/bitmaelum-suite/tools/internal"
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

	tools.VaultPassword = opts.Password
	v := tools.OpenVault()

	for i := range v.Accounts {
		var pow *proofofwork.ProofOfWork

		if len(v.Accounts[i].Pow.Data) == 0 {
			var err error
			addr, err := address.New(v.Accounts[i].Address)
			if err != nil {
				panic(err)
			}

			// proof of work is actually our hash address
			v.Accounts[i].Pow.Data = addr.Hash().String()
		}

		if v.Accounts[i].Pow.Bits >= opts.Bits && v.Accounts[i].Pow.IsValid() && !opts.Force {
			fmt.Printf("Account %s has %d bits\n", v.Accounts[i].Address, v.Accounts[i].Pow.Bits)
			continue
		}

		fmt.Printf("Working on %s\n", v.Accounts[i].Address)
		pow = proofofwork.New(opts.Bits, v.Accounts[i].Pow.Data, 0)
		pow.WorkMulticore()

		v.Accounts[i].Pow = *pow

		err := v.Save()
		if err != nil {
			panic(err)
		}
	}
}
