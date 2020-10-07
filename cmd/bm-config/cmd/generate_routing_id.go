package cmd

import (
	"fmt"
	"os"
	"strings"

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
		seed, r, err := config.GenerateRouting()
		if err != nil {
			logrus.Fatalf("Error while generating routing file: %v", err)
		}

		// Save routing
		err = config.SaveRouting(config.Server.Server.RoutingFile, r)
		if err != nil {
			logrus.Fatalf("Error while creating routing file: %v", err)
		}

		logrus.Printf("Generated routing file: %s", config.Server.Server.RoutingFile)

		fmt.Print(`
*****************************************************************************
!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORTANT!IMPORT
*****************************************************************************

We have generated a private key which allows you to control your server. If, 
for any reason, you lose this key, you will need to use the following words 
in order to recreate the key:

`)
		fmt.Print(wordWrap(seed, 78))
		fmt.Print(`

Write these words down and store them in a secure environment. They are the 
ONLY way to recover your private key in case you lose it.
`)
	},
}

func init() {
	rootCmd.AddCommand(initRoutingConfigCmd)
}

func wordWrap(s string, limit int) string {
	if strings.TrimSpace(s) == "" {
		return s
	}

	words := strings.Fields(strings.ToUpper(s))

	var result, line string
	for len(words) > 0 {
		if len(line)+len(words[0]) > limit {
			result += strings.TrimSpace(line) + "\n"
			line = ""
		}

		line = line + words[0] + " "
		words = words[1:]
	}
	if line != "" {
		result += strings.TrimSpace(line)
	}

	return result
}
