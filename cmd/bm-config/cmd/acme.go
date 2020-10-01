package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-config/internal/fileio"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-config/internal/letsencrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/spf13/cobra"
)

// letsEncryptCmd represents the lets encrypt  command
var letsEncryptCmd = &cobra.Command{
	Use:     "lets-encrypt",
	Aliases: []string{"le", "letsencrypt", "acme"},
	Short:   "Generates a certificate through LetsEncrypt",
	Long: `This command generates a signed certificate for your BitMaelum server through LetsEncrypt. 

Note that this command needs to setup a (temporary) HTTP web server on port 80, which means it 
needs to be run as root/administrator

It's possible to run the HTTP server on another port (with the --port/-p option), however, 
LetsEncrypt will ALWAYS connect to port 80. This flag is here so you can proxy your port 80 
to another (non-privileged) port.

This command will store LetsEncrypt account information in your acme path and store the 
certificate and key in the paths you defined in your config.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !config.Server.Acme.Enabled {
			fmt.Println("LetsEncrypt certificate generation is disabled. Check your configuration for more information.")
			os.Exit(1)
		}

		do(*staging, *port)
	},
}

func do(useStaging bool, httpPort string) {
	// 1.  Get or register ACME account
	fmt.Printf(" * Fetching account from your ACME directory...")
	le, err := letsencrypt.New(httpPort, useStaging, getAcmePath(useStaging))
	if err != nil {
		fmt.Printf("An error occurred: %s\n", err)
		os.Exit(1)
	}
	defer le.Cancel()

	if le.Account == nil {
		fmt.Println("not found")

		fmt.Printf(" * Registering a new account at LetsEncrypt...")
		err := le.RegisterAccount(config.Server.Acme.Email)
		if err != nil {
			fmt.Printf("An error occurred: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

		// Save into account directory. We don't care if it fails
		fmt.Printf(" * Saving account info in your ACME directory...")
		err = le.SaveAccount(getAcmePath(useStaging))
		if err != nil {
			fmt.Printf("An error occurred: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	} else {
		fmt.Println("found")
	}

	// 2.  Check if we are renewing an existing certificate or not
	fmt.Printf(" * Checking current certificate in \"%s\"...", config.Server.Server.CertFile)
	cert := le.LoadCertificate(config.Server.Server.CertFile)
	if cert == nil {
		fmt.Println("not found. That's ok")
	} else {
		fmt.Println("found")

		if cert.Subject.CommonName != config.Server.Acme.Domain {
			fmt.Printf(" * Domain found in certificate does not match configured name (want: %s, got: %s)\n", config.Server.Acme.Domain, cert.Subject.CommonName)
			os.Exit(1)
		}

		days, err := strconv.Atoi(config.Server.Acme.RenewBeforeDays)
		if err != nil {
			days = 30
		}

		if le.CheckRenewal(cert, days) {
			fmt.Println(" * Domain does not need to be renewed yet.")
			os.Exit(0)
		}
		fmt.Printf(" * Domain \"%s\" valid until %s\n", cert.Subject.CommonName, cert.NotAfter.Format(time.RFC822))
	}

	// 3.  Starting HTTP server for validation
	fmt.Printf(" * Starting up the HTTP server...")
	server := le.StartHTTPServer()
	defer func() {
		_ = server.Close()
	}()
	fmt.Println("up")

	// 4.  Submit an order for a certificate to be issued
	fmt.Printf(" * Creating an order at LetsEncrypt for our domain...")
	order, err := le.AuthorizeOrder(config.Server.Acme.Domain)
	if err != nil {
		fmt.Println("An error occurred while trying to register your domain at LetsEncrypt")
		os.Exit(1)
	}
	fmt.Println("ok")

	// 5.  Prove control of any identifiers requested in the certificate
	fmt.Printf(" * Waiting until LetsEncrypt calls our webserver...")
	order, err = le.Authorize(order)
	if err != nil {
		fmt.Printf("An error occurred while trying to authorize your domain at LetsEncrypt: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")

	// 6.  Finalize the order by submitting a CSR
	fmt.Printf(" * Verification complete. Asking for your certificate...")
	privKeyPem, certPem, err := le.FinalizeOrder(order, config.Server.Acme.Domain, nil)
	if err != nil {
		fmt.Printf("Error while fetching certificates: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")

	// 7. Save the files
	fmt.Println(" * Writing files:")
	err = fileio.SaveCertFiles(certPem, privKeyPem)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("\nAll done. Reload your BitMaelum server and you are good to go!\n")
}

func getAcmePath(useStaging bool) string {
	if useStaging {
		return filepath.Join(config.Server.Acme.Path, "staging")
	}
	return filepath.Join(config.Server.Acme.Path, "prod")
}

var (
	staging *bool
	port    *string
)

func init() {
	rootCmd.AddCommand(letsEncryptCmd)

	staging = letsEncryptCmd.Flags().BoolP("staging", "s", false, "Use Lets Encrypt staging environment")
	port = letsEncryptCmd.Flags().StringP("port", "p", "80", "TCP Port to start the webserver on")

	_ = letsEncryptCmd.MarkFlagRequired("email")
}
