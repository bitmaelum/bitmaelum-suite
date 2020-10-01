package cmd

import (
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initRoutingConfigCmd represents the initRoutingConfig command
var initRoutingConfigCmd = &cobra.Command{
	Use:   "generate-routing-id",
	Short: "Generates a routing ID and keypair",
	Long: `Before you can run the mail server, you will need a routing configuration file which uniquely identifies 
the server on the network.

This command creates a new routing file if one does not exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Server.Server.RoutingFile == "" {
			logrus.Fatalf("Routing file path is not found in your server configuration.")
		}

		// Check if file exist
		_, err := os.Stat(config.Server.Server.RoutingFile)
		if os.IsExist(err) {
			logrus.Fatalf("Routing file %s already exist. I will not overwrite this file.", config.Server.Server.RoutingFile)
		}

		// Generate new routing
		r, err := config.Generate()
		if err != nil {
			logrus.Fatalf("Error while generating routing file: %v", err)
		}

		// Save routing
		err = config.SaveRouting(config.Server.Server.RoutingFile, r)
		if err != nil {
			logrus.Fatalf("Error while creating routing file: %v", err)
		}

		logrus.Printf("Generated routing file: %s", config.Server.Server.RoutingFile)
	},
}

func init() {
	rootCmd.AddCommand(initRoutingConfigCmd)
}
