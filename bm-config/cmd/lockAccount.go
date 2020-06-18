package cmd

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/account/client"
	"github.com/bitmaelum/bitmaelum-server/core/password"
	"github.com/opentracing/opentracing-go/log"

	"github.com/spf13/cobra"
)

var lockAccountCmd = &cobra.Command{
	Use:   "lock-account",
	Short: "Lock your account with a password",
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := core.NewAddressFromString(*address)
		if err != nil {
			log.Error(err)
			return
		}

		if client.IsLocked(*addr) {
			fmt.Printf("This account is already locked.\n")
			return
		}

		err = client.LockAccount(*addr, password.AskPassword())
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Printf("This account is now locked.\n")
	},
}

var address *string

func init() {
	rootCmd.AddCommand(lockAccountCmd)

	address = lockAccountCmd.Flags().String("address", "","Address to lock")
	_ = lockAccountCmd.MarkFlagRequired("address")
}
