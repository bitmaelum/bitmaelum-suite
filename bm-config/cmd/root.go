package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bm-config",
	Short: "Configuration application for your mail server and client",
	Long: `This tool allows you to easily manage certain aspects of your mail server and client`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	allowRegistrationCmd.PersistentFlags().StringArrayP("config", "c", []string{}, "configuration file")
}
