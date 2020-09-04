package letsencrypt

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-config/internal/fileio"
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"golang.org/x/crypto/acme"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// LetsEncrypt structure holds everything for generating certificates through the LetsEncrypt ACME provider
type LetsEncrypt struct {
	AcmeClient *acme.Client
	Account    *acme.Account
	Key        crypto.Signer
	Ctx        context.Context
	Cancel     context.CancelFunc
	HTTPPort   string
}

// New creates a new LetsEncrypt instance and loads account info (if present)
func New(httpPort string, useStaging bool, p string) (*LetsEncrypt, error) {
	le := &LetsEncrypt{
		HTTPPort: httpPort,
	}

	// Use staging environment if needed
	var acmeDir = acme.LetsEncryptURL
	if useStaging {
		acmeDir = "https://acme-staging-v02.api.letsencrypt.org/directory"
	}

	le.AcmeClient = &acme.Client{
		DirectoryURL: acmeDir,
		UserAgent:    "BitMaelum v" + internal.Version.String(),
	}

	// Create context
	le.Ctx, le.Cancel = context.WithTimeout(context.Background(), 2*time.Minute)

	// Load account information
	key, account, err := getAccountFromAcmeDir(p)
	if err == nil {
		le.Key = key
		le.Account = account
		le.AcmeClient.Key = key

		return le, nil
	}

	return le, nil
}

// RegisterAccount will register a new account with LetsEncrypt based on the given email. It automatically generates
// a new private key which is used for communication. This key must be set in the AcmeClient in order to communicate
// with LetsEncrypt.
func (le *LetsEncrypt) RegisterAccount(email string) error {
	// Create public/private key pair.
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	// Needed to set here so acmeClient.register() will work
	le.AcmeClient.Key = privKey

	account := &acme.Account{
		Contact: []string{"mailto:" + email},
	}
	account, err = le.AcmeClient.Register(le.Ctx, account, acme.AcceptTOS)
	if err != nil {
		return err
	}

	// Set the rest of the settings
	le.Key = privKey
	le.Account = account

	return nil
}

// CheckRenewal will check if the renewal date of the given certificate is met. This is the expiry-date of the certificate
// minus a number of days (default 30).
func (le *LetsEncrypt) CheckRenewal(cert *x509.Certificate, days int) bool {
	d := time.Duration(days) * 24 * time.Hour
	renewAfter := cert.NotAfter.Add(-1 * d).Unix()

	if renewAfter < time.Now().Unix() {
		return false
	}

	return true
}

// LoadCertificate loads a certificate from the given path or returns nil when no certificate is found.
func (le *LetsEncrypt) LoadCertificate(p string) *x509.Certificate {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil
	}

	// Decode FIRST certificate found in the certfile
	var certDERBlock *pem.Block
	certDERBlock, _ = pem.Decode(data)
	cert, err := x509.ParseCertificates(certDERBlock.Bytes)
	if err != nil || len(cert) == 0 {
		return nil
	}

	return cert[0]
}

// FinalizeOrder will ask LetsEncrypt for the actual certificate by issusing a certificate order. Note that this can
// only be done with a valid order. When privCertKey is nil, it will generate a new random RSA 2048 key. When we are
// renewing a certificate you might want to keep the same key (or not).
func (le *LetsEncrypt) FinalizeOrder(order *acme.Order, domain string, privCertKey crypto.Signer) (privKey string, certificate string, err error) {
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
	certs, _, err := le.AcmeClient.CreateOrderCert(le.Ctx, order.FinalizeURL, csr, false)
	if err != nil {
		return "", "", err
	}

	// Convert all certificates to PEM format
	for _, cert := range certs {
		c, err := x509.ParseCertificate(cert)
		if err != nil {
			return "", "", err
		}

		var b bytes.Buffer
		err = pem.Encode(&b, &pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})
		if err != nil {
			return "", "", err
		}

		certificate += b.String()
	}

	// Convert private and public key to PEM
	privKey, err = privKeyToPEM(privCertKey)
	if err != nil {
		return "", "", err
	}

	return
}

// Authorize will get an unauthorized order and tries to authorize it. This is done by asking LetsEncrypt to validate
// the order. This is done by letsencrypt calling our HTTP server on which we respond with a special file.
func (le *LetsEncrypt) Authorize(order *acme.Order) (*acme.Order, error) {
	for _, authzURL := range order.AuthzURLs {
		auth, err := le.AcmeClient.GetAuthorization(le.Ctx, authzURL)
		if err != nil {
			continue
		}

		for _, challenge := range auth.Challenges {
			if challenge.Type == "http-01" {
				order, err := le.acceptChallenge(challenge, order)
				if err != nil {
					continue
				}

				return order, nil
			}
		}
	}

	return nil, errors.New("cannot find a authorization/challenge")
}

// acceptChallenge will accept a given challenge (http-01) for the given order and waits until LetsEncrypt has validated
// the order. It will then return a validated order which we can use for requesting the actual certificate through
// FinalizeOrder.
func (le *LetsEncrypt) acceptChallenge(challenge *acme.Challenge, order *acme.Order) (*acme.Order, error) {
	// Do this challenge
	_, err := le.AcmeClient.Accept(le.Ctx, challenge)
	if err != nil {
		return nil, err
	}

	// Wait until order is completed
	order, err = le.AcmeClient.WaitOrder(le.Ctx, order.URI)
	if err != nil {
		return nil, err
	}

	if order.Status != acme.StatusReady {
		return nil, errors.New("validation seems to have failed")
	}

	return order, nil
}

// SaveAccount will save the given ACME account and matching private Key into the Acme directory. This allows us to
// use the same account for requesting new certificates later on.
func (le *LetsEncrypt) SaveAccount(dir string) error {
	if le.Account == nil {
		return errors.New("account not loaded or registered yet")
	}

	// Make sure directory exists before writing
	_ = os.MkdirAll(dir, 0777)

	// @TODO: What happens if we mix staging and production.. or different accounts?
	err := fileio.SaveFile(filepath.Join(dir, "account.json"), le.Account)
	if err != nil {
		return err
	}

	// Convert key to PEM format
	s, err := privKeyToPEM(le.Key)
	if err != nil {
		return err
	}

	return fileio.SaveFile(filepath.Join(dir, "key.json"), s)
}

// StartHTTPServer will start the HTTP server in the background which is called by LetsEncrypt for validating orders.
func (le *LetsEncrypt) StartHTTPServer() *http.Server {
	// Setup simple web-server that listens to tokens and returns back.. things
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.URL.String(), "/.well-known/acme-challenge/")
		response, err := le.AcmeClient.HTTP01ChallengeResponse(token)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	})

	server := &http.Server{Addr: ":" + le.HTTPPort}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("We cannot start a HTTP server on port %s on your machine.\n", le.HTTPPort)
			if err.(*net.OpError).Err.(*os.SyscallError).Err == syscall.EADDRINUSE {
				fmt.Printf("It seems there is already a server running on port %s.\n", le.HTTPPort)
			} else if err.(*net.OpError).Err.(*os.SyscallError).Err == syscall.EACCES {
				fmt.Printf("It seems you are not running as root/Administrator, which is needed to use port %s.\n", le.HTTPPort)
			} else {
				fmt.Printf("An unknown error has occurred: %s\n", err)
			}
			os.Exit(1)
		}
	}()

	return server
}

// AuthorizeOrder will create a new (unvalidated) order for the given domain. It must first be validated through the
// Authorize method before we can call FinalizeOrder to fetch our certificate.
func (le *LetsEncrypt) AuthorizeOrder(domain string) (*acme.Order, error) {
	return le.AcmeClient.AuthorizeOrder(le.Ctx, acme.DomainIDs(domain))
}

// getAccountFromAcmeDir will return saved account information and private key from our acme directory.
func getAccountFromAcmeDir(dir string) (crypto.Signer, *acme.Account, error) {
	acc := &acme.Account{}
	err := fileio.LoadFile(filepath.Join(dir, "account.json"), &acc)
	if err != nil {
		return nil, nil, err
	}

	// can't marshal ecdsa.PrivKey to JSON (go 1.16+ probably)
	var data string
	err = fileio.LoadFile(filepath.Join(dir, "key.json"), &data)
	if err != nil {
		return nil, nil, err
	}
	privKey, err := bmcrypto.NewPrivKey(data)
	if err != nil {
		return nil, nil, err
	}

	return privKey.K.(crypto.Signer), acc, nil
}

// PrivKeyToPEM Convert a private key into PKCS8/PEM format
func privKeyToPEM(key interface{}) (string, error) {
	privBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = pem.Encode(&b, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	return b.String(), err
}
