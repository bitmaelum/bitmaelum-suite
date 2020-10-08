package cmd

import (
	"fmt"

	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var setDefaultCmd = &cobra.Command{
	Use:     "set-default",
	Aliases: []string{"default"},
	Short:   "Set default account to send from",
	Long:    `When you don't specify a from address, bm-client (ano other tools) will automatically use this address.`,
	Run: func(cmd *cobra.Command, args []string) {
		v := OpenVault()

		fromAddr, err := address.NewAddress(*sdFrom)
		if err != nil {
			logrus.Fatal(err)
		}
		if !v.HasAccount(*fromAddr) {
			logrus.Fatal("Address not found in vault")
		}

		for i := range v.Store.Accounts {
			v.Store.Accounts[i].Default = false

			if v.Store.Accounts[i].Address == *sdFrom {
				v.Store.Accounts[i].Default = true
			}
		}

		err = v.WriteToDisk()
		if err != nil {
			logrus.Fatalf("error while saving vault: %s\n", err)
		}

		fmt.Printf("Default address set to '%s'\n", *sdFrom)
	},
}

var sdFrom *string

func init() {
	rootCmd.AddCommand(setDefaultCmd)

	sdFrom = setDefaultCmd.Flags().String("address", "", "Default address to set")

	_ = setDefaultCmd.MarkFlagRequired("address")
}
