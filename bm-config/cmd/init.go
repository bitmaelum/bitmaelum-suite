package cmd

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var (
	clientConfigPath string = "./client.config.yml"
	serverConfigPath string = "./server.config.yml"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates default server and client configurations",
	Long: `Before you can run the mailserver or client, you will need a configuration file which you need to adjust 
to your own needs.

This command creates default templates that you can use as a starting point.`,
	Run: func(cmd *cobra.Command, args []string) {
		createFile(clientConfigPath, config.GenerateClientConfig)
		createFile(serverConfigPath, config.GenerateServerConfig)
	},
}

func createFile(path string, configTemplate func(w io.Writer) error) {
	_, err := os.Stat(path)
	if err != nil {
		// File exists
		fmt.Println(path + " already exists. Skipping.")
		return
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	err = configTemplate(f)
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
}
