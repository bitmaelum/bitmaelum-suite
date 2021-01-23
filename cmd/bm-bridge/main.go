// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	imapgw "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal/imap/backend"
	smtpgw "github.com/bitmaelum/bitmaelum-suite/cmd/bm-bridge/internal/smtp/backend"
	bminternal "github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/emersion/go-imap/server"
	"github.com/emersion/go-smtp"
)

type options struct {
	ImapHost string `long:"imaphost" description:"Host:Port to imap server from" required:"true"`
	SmtpHost string `long:"smtphost" description:"Host:Port to smtp server from" required:"true"`
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
	Vault    string `long:"vault" description:"Custom vault file" default:""`
	Version  bool   `short:"v" long:"version" description:"Display version information"`
}

// Opts are the options set through the command line
var opts options

func main() {
	rand.Seed(time.Now().UnixNano())

	bminternal.ParseOptions(&opts)
	if opts.Version {
		bminternal.WriteVersionInfo("BitMaelum IMAP", os.Stdout)
		fmt.Println()
		os.Exit(0)
	}

	config.LoadClientConfig(opts.Config)

	// Set default vault info if set in config
	vault.VaultPassword = opts.Password
	vault.VaultPath = config.Client.Vault.Path
	if opts.Vault != "" {
		vault.VaultPath = opts.Vault
	}

	v := vault.OpenDefaultVault()

	fmt.Println(" * BM-BRIDGE *")

	// setup context so we can easily stop all components of the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Wait for signals and cancel context
	go setupSignals(cancel)

	go startImapServer(v, cancel)
	go startSmtpServer(v, cancel)

	// @TODO: SMTP Server

	<-ctx.Done()

}

func setupSignals(cancel context.CancelFunc) {
	go func() {
		// Capture INT and TERM signals
		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

		for {
			s := <-sigChannel
			fmt.Printf("Signal %s received. Terminating.\n", s)
			cancel()
		}
	}()
}

func startImapServer(v *vault.Vault, cancel context.CancelFunc) {
	be := imapgw.New(v)
	s := server.New(be)
	s.Addr = opts.ImapHost
	s.AllowInsecureAuth = true // We should run on TLS
	s.Debug = os.Stdout

	fmt.Println(" [*] IMAP starting")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal("error while accepting connection: " + err.Error() + "\n")
		cancel()
	}
}

func startSmtpServer(v *vault.Vault, cancel context.CancelFunc) {
	be := smtpgw.New(v)
	s := smtp.NewServer(be)
	s.Addr = opts.SmtpHost
	s.Domain = "bitmaelum.network"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 1
	s.AllowInsecureAuth = true
	s.Debug = os.Stdout

	fmt.Println(" [*] SMTP starting")
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("error while accepting connection: " + err.Error() + "\n")
		cancel()
	}
}
