package cmd

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/account/client"
	"github.com/bitmaelum/bitmaelum-server/core/password"
	"github.com/opentracing/opentracing-go/log"

	"github.com/spf13/cobra"
)

var unlockAccountCmd = &cobra.Command{
	Use:   "unlock-account",
	Short: "Unlock an account. No password is needed to send/receive mail",
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := core.NewAddressFromString(*address)
		if err != nil {
			log.Error(err)
			return
		}

		if ! client.IsLocked(*addr) {
			fmt.Printf("This account is already unlocked.\n")
			return
		}

		err = client.UnlockAccount(*addr, password.AskPassword())
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Printf("This account is now unlocked.\n")
	},
}

func init() {
	rootCmd.AddCommand(unlockAccountCmd)

	address = unlockAccountCmd.Flags().String("address", "","Address to unlock")
	_ = unlockAccountCmd.MarkFlagRequired("address")
}
