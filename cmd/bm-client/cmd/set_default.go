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
		vault := OpenVault()

		fromAddr, err := address.New(*sdFrom)
		if err != nil {
			logrus.Fatal(err)
		}
		if ! vault.HasAccount(*fromAddr) {
			logrus.Fatal("Address not found in vault")
		}

		for i := range vault.Accounts {
			vault.Accounts[i].Default = false

			if vault.Accounts[i].Address == *sdFrom {
				vault.Accounts[i].Default = true
			}
		}

		err = vault.Save()
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
