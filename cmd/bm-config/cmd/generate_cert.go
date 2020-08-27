package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"time"
)

// generateCertCmd represents the generateCert command
var generateCertCmd = &cobra.Command{
	Use:   "generate-cert",
	Short: "Generates a self-signed server certificate",
	Long: `This command generates a self-signed certificate for your BitMaelum server.

Note that self-signed servers can be used but not all mail-servers will accept 
self-signed certificates.`,
	Run: func(cmd *cobra.Command, args []string) {
		domainName, err := cmd.Flags().GetString("domain")
		if err != nil {
			logrus.Fatalf("Cannot read domain: %v", err)
		}
		generateCert(domainName)
	},
}

// Generate a self-signed certificate
// Taken mostly from https://golang.org/src/crypto/tls/generate_cert.go
func generateCert(domain string) {
	// Generate x509 template
	var notBefore = time.Now()
	var notAfter = notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logrus.Fatalf("Cannot generate serial number: %v ", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: domain,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	// Generate Private/Public RSA key
	fmt.Println("Generating 2048 bits keypair...")

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logrus.Fatalf("Error while generating keypair: %v", err)
	}

	// Create certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		logrus.Fatalf("Failed to generate certificate: %v", err)
	}

	// Write certificate to file
	fmt.Println("Writing server.cert file")

	certOut, err := os.Create("server.cert")
	if err != nil {
		logrus.Fatalf("Error while opening server.cert: %v", err)
	}
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		logrus.Fatalf("Failed to write to server.cert: %v", err)
	}
	err = certOut.Close()
	if err != nil {
		logrus.Fatalf("ERror while closing server.cert: %v", err)
	}

	// Write key to file
	fmt.Println("Writing server.key file")

	keyOut, err := os.Create("server.key")
	if err != nil {
		logrus.Fatalf("Error while opening server.key: %v", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		logrus.Fatalf("Error while marshalling key: %v", err)
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		logrus.Fatalf("Failed to write to server.key: %v", err)
	}
	err = keyOut.Close()
	if err != nil {
		logrus.Fatalf("Error while closing server.key: %v", err)
	}

	fmt.Println("All done.")
}

func init() {
	rootCmd.AddCommand(generateCertCmd)

	generateCertCmd.Flags().StringP("domain", "d", "localhost", "The common name / domain name you want in your certificate.")
}
