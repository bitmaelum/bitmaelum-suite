package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var setCmd = &cobra.Command{
	Use:     "set",
	Short:   "Create or update a specific setting in your vault",
	Long:    `Your vault accounts can have additional settings. With this command you can easily manage these.`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := OpenVault()
		info := GetAccountOrDefault(vault, *sAddress)

		var msg string

		k := strings.ToLower(*sKey)
		switch k {
		case "name":
			info.Name = *sValue
			msg = fmt.Sprintf("Updated setting %s to %s", "name", *sValue)
		case "organisation":
			info.Organisation = *sValue
			msg = fmt.Sprintf("Updated setting %s to %s", "organisation", *sValue)
		default:
			if *sValue == "" {
				if info.Settings != nil {
					delete(info.Settings, k)
				}
				msg = fmt.Sprintf("Removed setting %s", k)
			} else {
				if info.Settings == nil {
					info.Settings = make(map[string]string)
				}
				info.Settings[k] = *sValue
				msg = fmt.Sprintf("Updated setting %s to %s", k, *sValue)
			}
		}

		err := vault.Save()
		if err != nil {
			logrus.Fatalf("error while saving vault: %s\n", err)
		}

		logrus.Printf("%s\n", msg)
	},
}

var (
	sAddress *string
	sKey     *string
	sValue   *string
)

func init() {
	rootCmd.AddCommand(setCmd)

	sAddress = setCmd.Flags().String("address", "", "Default address to set")
	sKey = setCmd.Flags().String("key", "", "Key to set")
	sValue = setCmd.Flags().String("value", "", "Value to set")
}
