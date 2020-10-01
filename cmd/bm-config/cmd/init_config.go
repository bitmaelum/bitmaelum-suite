package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// initConfigCmd represents the initConfig command
var initConfigCmd = &cobra.Command{
	Use:   "init-config",
	Short: "Creates default server and client configurations",
	Long: `Before you can run the mailserver or client, you will need a configuration file which you need to adjust 
to your own needs.

This command creates default templates that you can use as a starting point.`,
	Run: func(cmd *cobra.Command, args []string) {
		genC, _ := cmd.Flags().GetBool("client")
		genS, _ := cmd.Flags().GetBool("server")

		// If not client or server selected, generate them both
		genBoth := !genC && !genS

		if genBoth || genC {
			createFile("./"+config.ClientConfigFile, config.GenerateClientConfig)
			fmt.Println("Generated client configuration file")
		}
		if genBoth || genS {
			createFile("./"+config.ServerConfigFile, config.GenerateServerConfig)
			fmt.Println("Generated server configuration file")
		}
	},
}

func createFile(path string, configTemplate func(w io.Writer) error) {
	f, err := os.Create(path)
	if err != nil {
		logrus.Fatalf("Error while creating file: %v", err)
	}

	err = configTemplate(f)
	if err != nil {
		logrus.Fatalf("Error while creating file: %v", err)
	}

	err = f.Close()
	if err != nil {
		logrus.Fatalf("Error while closing file: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(initConfigCmd)

	initConfigCmd.Flags().Bool("client", false, "Generate only the client configuration")
	initConfigCmd.Flags().Bool("server", false, "Generate only the server configuration")
}
