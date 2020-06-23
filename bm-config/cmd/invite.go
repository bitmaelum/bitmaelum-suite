package cmd

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/spf13/cobra"
	"time"
)

// inviteCmd represents the invite command
var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Invite a new user onto your server",
	Long: `This command will generate an invitation token that must be used for registering an account on your 
server. Only the specified address can register the account`,
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

		fmt.Printf("'%s' is allowed to register on our server in the next %d days.\n", addr.String(), d)
		fmt.Printf("The invitation token is: %s\n", token)
	},
}

func init() {
	rootCmd.AddCommand(inviteCmd)

	inviteCmd.Flags().String("address", "", "Address to register")
	inviteCmd.Flags().Int("days", 30, "Days allowed for registration")

	_ = inviteCmd.MarkFlagRequired("address")
}
