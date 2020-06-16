package cmd

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/spf13/cobra"
	"time"
)

// allowRegistrationCmd represents the allowRegistration command
var allowRegistrationCmd = &cobra.Command{

	Use:   "allow-registration",
	Short: "Allows registration of given address",
	Long: `When running a mailserver, it's nice to limit the number of users that can create addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := cmd.Flags().GetString("address")
		d, _ := cmd.Flags().GetInt("days")

		addr, err := core.NewAddressFromString(s)
		if err != nil {
			fmt.Printf("incorrect address specified")
			return
		}

		is := container.GetInviteService()
		token, err := is.GetInvite(addr.Hash())
		if err == nil {
			fmt.Printf("'%s' already allowed to register with token: %s\n", addr.String(), token)
			return
		}

		token, err = is.CreateInvite(addr.Hash(), time.Duration(d) * 24 * time.Hour)
		if err != nil {
			fmt.Printf("error while inviting address")
		}

		fmt.Printf("'%s' is allowed to register on our server in the next %d days. The token is: %s\n", addr.String(), d, token)
	},
}

func init() {
	rootCmd.AddCommand(allowRegistrationCmd)

	allowRegistrationCmd.Flags().String("address", "", "Address to register")
	allowRegistrationCmd.Flags().Int("days", 30, "Days allowed for registration")
}
