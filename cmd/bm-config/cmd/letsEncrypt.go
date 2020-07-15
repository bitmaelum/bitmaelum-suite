package cmd

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	ctx        context.Context
	acmeClient *acme.Client
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
		if config.Server.Acme.Enabled == false {
			fmt.Println("LetsEncrypt certificate generation is disabled. Check your configuration for more information.")
			os.Exit(1)
		}

		registerLetsEncrypt(*staging, *port)
	},
}

func registerLetsEncrypt(useStaging bool, httpPort string) {
	var (
		err    error
		domain = config.Server.Server.Name
	)

	// Check if we are renewing an existing certificate or not
	var cert, renewalMode = checkRenewal()
	if cert != nil && renewalMode == false {
		fmt.Println(" * Domain does not need to be renewed yet.")
		os.Exit(0)
	}

	// Create context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Create client
	fmt.Printf(" * Trying to setup the connection to LetsEncrypt...")
	acmeClient, err = getClient(useStaging)
	if err != nil {
		fmt.Println("Cannot connect to LetsEncrypt")
		os.Exit(1)
	}
	fmt.Println("ok")

	// Starting HTTP server for validation
	fmt.Printf(" * Starting up the HTTP server...")
	server := startHTTPServer(httpPort)
	defer func() {
		_ = server.Close()
	}()
	fmt.Println("up")

	// Register or fetch account and private key, and set it in our acme client
	accountPrivKey, _, err := getOrRegisterAccount(config.Server.Acme.Email)
	acmeClient.Key = accountPrivKey

	// 1.  Submit an order for a certificate to be issued
	fmt.Printf(" * Creating an order at LetsEncrypt for our domain...")
	order, err := acmeClient.AuthorizeOrder(ctx, acme.DomainIDs(domain))
	if err != nil {
		fmt.Println("An error occurred while trying to register your domain at LetsEncrypt")
		os.Exit(1)
	}
	fmt.Println("ok")

	// 2.  Prove control of any identifiers requested in the certificate
	fmt.Printf(" * Waiting until LetsEncrypt calls our webserver...")
	order, err = authorize(order)
	if err != nil {
		fmt.Printf("An error occurred while trying to authorize your domain at LetsEncrypt: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")

	// 3.  Finalize the order by submitting a CSR
	fmt.Printf(" * Verification complete. Asking for your certificate...")
	privKeyPem, certPem, err := finalizeOrder(order, domain, nil)
	if err != nil {
		fmt.Printf("Error while fetching certificates: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("ok")

	// 4. Save the files
	fmt.Println(" * Writing files:")
	err = saveCertFiles(certPem, privKeyPem)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("\nAll done. Reload your BitMaelum server and you are good to go!\n")
}

func saveCertFiles(certPem string, keyPem string) error {
	// Find the latest suffix we can use that both certfile and keyfile do not use.
	var suffix = 1
	for {
		p, _ := homedir.Expand(fmt.Sprintf("%s.%03d", config.Server.Server.CertFile, suffix))
		_, err1 := os.Stat(p)
		p, _ = homedir.Expand(fmt.Sprintf("%s.%03d", config.Server.Server.KeyFile, suffix))
		_, err2 := os.Stat(p)
		if err1 == nil || err2 == nil {
			suffix++
			continue
		}
		break
	}

	newPath, _ := homedir.Expand(fmt.Sprintf("%s.%03d", config.Server.Server.CertFile, suffix))
	oldPath, _ := homedir.Expand(config.Server.Server.CertFile)
	fmt.Printf("   - moving old cert file to %s: ", newPath)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	newPath, _ = homedir.Expand(fmt.Sprintf("%s.%03d", config.Server.Server.KeyFile, suffix))
	oldPath, _ = homedir.Expand(config.Server.Server.KeyFile)
	fmt.Printf("   - moving old key file to %s: ", newPath)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}
	fmt.Println("ok")


	fmt.Printf("   - Writing new cert file %s: ", config.Server.Server.CertFile)
	newPath, _ = homedir.Expand(config.Server.Server.CertFile)
	err = ioutil.WriteFile(newPath, []byte(certPem), 0600)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	fmt.Printf("   - Writing new key file %s: ", config.Server.Server.CertFile)
	newPath, _ = homedir.Expand(config.Server.Server.KeyFile)
	err = ioutil.WriteFile(newPath, []byte(keyPem), 0600)
	if err != nil {
		return err
	}
	fmt.Println("ok")

	return nil
}

func getOrRegisterAccount(email string) (crypto.Signer, *acme.Account, error) {
	// 0. Load or register account
	fmt.Printf(" * Fetching account from your ACME directory...")
	accountPrivKey, acc, err := getAccountFromAcmeDir()
	if err == nil {
		fmt.Println("found")

		return accountPrivKey, acc, nil
	}
	fmt.Println("not found.")

	// Saved acme account not found. Register a new one..
	fmt.Printf(" * Registering your email at LetsEncrypt...")
	accountPrivKey, acc, err = registerAccount(email)
	if err != nil {
		fmt.Println("An error occurred while trying to register your email address. Maybe it's incorrect?")
		os.Exit(1)
	}
	fmt.Println("ok")

	// Save into account directory. We don't care if it fails
	fmt.Printf(" * Saving account in your ACME directory...")
	_ = saveAccount(acc, accountPrivKey)
	fmt.Println("ok")

	return accountPrivKey, acc, nil
}

func checkRenewal() (*x509.Certificate, bool) {
	fmt.Printf(" * Checking current certificate in \"%s\"...", config.Server.Server.CertFile)
	cert := fetchCertificate()
	if cert == nil {
		fmt.Println("not found.")
		return nil, false
	}

	fmt.Println("found")

	host, _, _ := net.SplitHostPort(config.Server.Server.Name)
	if host == "" {
		host = config.Server.Server.Name
	}

	if cert.Subject.CommonName != host {
		fmt.Printf(" * Domain found in certificate does not match configured name (want: %s, got: %s)\n", host, cert.Subject.CommonName)
		os.Exit(1)
	}

	fmt.Printf(" * Domain \"%s\" valid until %s\n", cert.Subject.CommonName, cert.NotAfter.Format(time.RFC822))
	days, err := strconv.Atoi(config.Server.Acme.RenewBeforeDays)
	if err != nil {
		days = 30
	}
	if cert.NotAfter.Unix() > time.Now().Add(time.Duration(days*24)*time.Hour).Unix() {
		return cert, false
	}

	return cert, true
}

func fetchCertificate() *x509.Certificate {
	p, err := homedir.Expand(config.Server.Server.CertFile)
	if err != nil {
		return nil
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil
	}

	// Decode FIRST certificate found in the certfile
	var certDERBlock *pem.Block
	certDERBlock, data = pem.Decode(data)
	cert, err := x509.ParseCertificates(certDERBlock.Bytes)
	if err != nil || len(cert) == 0 {
		return nil
	}

	return cert[0]
}

func finalizeOrder(order *acme.Order, domain string, privCertKey crypto.Signer) (privKey string, certificate string, err error) {
	// Generate a new private key if not supplied
	if privCertKey == nil {
		privCertKey, err = rsa.GenerateKey(rand.Reader, 2048)
	}

	// Create CSR
	req := &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: domain},
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, req, privCertKey)
	if err != nil {
		return "", "", err
	}

	// Let LetsEncrypt sign the CSR and return the certificate
	certs, _, err := acmeClient.CreateOrderCert(ctx, order.FinalizeURL, csr, false)
	if err != nil {
		return "", "", err
	}

	// Convert all certificates to PEM format
	for _, cert := range certs {
		c, err := x509.ParseCertificate(cert)
		if err != nil {
			return "", "", err
		}

		certPem, err := encrypt.CertToPEM(*c)
		certificate += certPem
	}

	// Convert private and public key to PEM
	privKey, err = encrypt.PrivKeyToPEM(privCertKey)
	if err != nil {
		return "", "", err
	}

	return
}

func authorize(order *acme.Order) (*acme.Order, error) {
	for _, authzUrl := range order.AuthzURLs {
		auth, err := acmeClient.GetAuthorization(ctx, authzUrl)
		if err != nil {
			continue
		}

		for _, challenge := range auth.Challenges {
			if challenge.Type == "http-01" {
				order, err := acceptChallenge(ctx, challenge, order)
				if err != nil {
					continue
				}

				return order, nil
			}
		}
	}

	return nil, errors.New("cannot find a authorization/challenge")
}

func acceptChallenge(ctx context.Context, challenge *acme.Challenge, order *acme.Order) (*acme.Order, error) {
	// Do this challenge
	_, err := acmeClient.Accept(ctx, challenge)
	if err != nil {
		return nil, err
	}

	// Wait until order is completed
	order, err = acmeClient.WaitOrder(ctx, order.URI)
	if err != nil {
		return nil, err
	}

	if order.Status != acme.StatusReady {
		return nil, errors.New("validation seems to have failed")
	}

	return order, nil
}

func registerAccount(email string) (crypto.Signer, *acme.Account, error) {
	// Create public/private key pair.
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	// Needed to set here so acmeClient.register() will work
	acmeClient.Key = key

	acc := &acme.Account{
		Contact: []string{"mailto:" + email},
	}
	acc, err = acmeClient.Register(ctx, acc, acme.AcceptTOS)
	if err != nil {
		return nil, nil, err
	}

	return key, acc, nil
}

func saveAccount(acc *acme.Account, key crypto.Signer) error {
	// @TODO: What happens if we mix staging and production.. or different accounts?
	err := saveFile(config.Server.Acme.Path+"/account.json", acc)
	if err != nil {
		return err
	}

	// Convert key to PEM format
	s, err := encrypt.PrivKeyToPEM(key)
	if err != nil {
		return err
	}

	return saveFile(config.Server.Acme.Path+"/key.json", s)
}

func saveFile(p string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	p, err = homedir.Expand(p)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Dir(p), 755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, data, 0600)
}

func getClient(useStaging bool) (*acme.Client, error) {
	// Use staging environment if needed
	var acmeDir = acme.LetsEncryptURL
	if useStaging {
		acmeDir = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	client := &acme.Client{
		DirectoryURL: acmeDir,
		UserAgent:    "BitMaelum v" + internal.Version.String(),
	}

	return client, nil
}

func startHTTPServer(port string) *http.Server {
	// Setup simple web-server that listens to tokens and returns back.. things
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(".")
		// fmt.Printf("   - Incoming request: '%s'\n", r.URL.String())

		token := strings.TrimPrefix(r.URL.String(), "/.well-known/acme-challenge/")
		// fmt.Printf("   - Found token: '%s'\n", token)
		response, err := acmeClient.HTTP01ChallengeResponse(token)
		// fmt.Printf("   - Sending response: '%s'\n", response)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	})

	server := &http.Server{Addr: ":" + port}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("We cannot start a HTTP server on port %s on your machine.\n", port)
			if err.(*net.OpError).Err.(*os.SyscallError).Err == syscall.EADDRINUSE {
				fmt.Printf("It seems there is already a server running on port %s.\n", port)
			} else if err.(*net.OpError).Err.(*os.SyscallError).Err == syscall.EACCES {
				fmt.Printf("It seems you are not running as root/Administrator, which is needed to use port %s.\n", port)
			} else {
				fmt.Printf("An unknown error has occurred: %s\n", err)
			}
			os.Exit(1)
		}
	}()

	return server
}

func getAccountFromAcmeDir() (crypto.Signer, *acme.Account, error) {
	acc := &acme.Account{}
	err := loadFile(config.Server.Acme.Path+"/account.json", &acc)
	if err != nil {
		return nil, nil, err
	}

	// can't marshal ecdsa.PrivKey to JSON (go 1.16+ probably)
	var data string
	err = loadFile(config.Server.Acme.Path+"/key.json", &data)
	if err != nil {
		return nil, nil, err
	}
	key, err := encrypt.PEMToPrivKey([]byte(data))
	if err != nil {
		return nil, nil, err
	}

	return key.(crypto.Signer), acc, nil
}

func loadFile(p string, v interface{}) error {
	p, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
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
