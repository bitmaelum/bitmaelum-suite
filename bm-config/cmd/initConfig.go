package cmd

import (
	"github.com/bitmaelum/bitmaelum-server/core/config"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
)

var (
	clientConfigPath string = "./client.config.yml"
	serverConfigPath string = "./server.config.yml"
)

// initConfigCmd represents the initConfig command
var initConfigCmd = &cobra.Command{
	Use:   "init-config",
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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Fatalf("Error while creating file: %v", err)
	}

	err = configTemplate(f)
	if err != nil {
		log.Fatalf("Error while creating file: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("Error while closing file: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(initConfigCmd)
}
