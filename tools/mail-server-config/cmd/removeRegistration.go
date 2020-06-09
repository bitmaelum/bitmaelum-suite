package cmd

import (
	"fmt"
	"github.com/jaytaph/mailv2/core"
	"github.com/jaytaph/mailv2/core/container"
	"github.com/spf13/cobra"
)

// removeRegistrationCmd represents the removeRegistration command
var removeRegistrationCmd = &cobra.Command{
	Use:   "remove-registration",
	Short: "Removes the registration of given address",
	Long: `When running a mailserver, it's nice to limit the number of users that can create addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := cmd.Flags().GetString("address")

		addr, err := core.NewAddressFromString(s)
		if err != nil {
			fmt.Printf("incorrect address specified")
			return
		}

		is := container.GetInviteService()
		err = is.RemoveInvite(addr.Hash())
		if err != nil {
			fmt.Printf("error while removing address")
		}

		fmt.Printf("'%s' has been removed.\n", addr.String())
	},
}

func init() {
	rootCmd.AddCommand(removeRegistrationCmd)

	removeRegistrationCmd.Flags().String("address", "", "Address to remove")
}
