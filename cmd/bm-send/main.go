// Copyright (c) 2020 BitMaelum Authors
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
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/jessevdk/go-flags"
)

type options struct {
	PrivateKey  string   `short:"p" long:"private_key" description:"Private key" required:"true" env:"BITMAELUM_SEND_PRIVATE_KEY"`
	Subject     string   `short:"s" long:"subject" description:"Subject" required:"true"`
	From        string   `short:"f" long:"from" description:"Sender" required:"true" env:"BITMAELUM_SEND_FROM"`
	To          string   `short:"t" long:"to" description:"Recipient" required:"true"`
	Message     string   `short:"m" long:"message" description:"Default message"`
	Blocks      []string `short:"b" long:"block" description:"Body block"`
	Attachments []string `short:"a" long:"attachment" description:"Attachment"`
	Resolver    string   `short:"r" long:"resolver" description:"Resolver" env:"BITMAELUM_SEND_RESOLVER_URL" default:"resolver.bitmaelum.com"`
}

var opts options

var (
	fromAddr *address.Address
	toAddr   *address.Address
	privKey  *bmcrypto.PrivKey
)

func apiErrorFunc(_ *http.Request, resp *http.Response) {
	// Read body
	b, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	err := api.GetErrorFromResponse(b)
	if err != nil {
		fmt.Println("error: ", err.Error())
		os.Exit(1)
	}

	// Whoops.. not an error. Let's pretend nothing happened and create a new buffer so we can read the body again
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
}


func main() {
	rand.Seed(time.Now().UnixNano())

	parseFlags()

	if opts.Message != "" && len(opts.Blocks) > 0 {
		fmt.Println("either you send a message, or one or more blocks, but not both")
		os.Exit(1)
	}

	// Set default message block if a message is specified
	if opts.Message != "" {
		opts.Blocks = append(opts.Blocks, "default:"+opts.Message)
	}

	// Set resolve settings, as we don't use a client configuration file
	config.Client.Resolver.Remote.Enabled = true
	config.Client.Resolver.Remote.URL = opts.Resolver

	// Fetch both sender and recipient info
	svc := container.Instance.GetResolveService()
	senderInfo, err := svc.ResolveAddress(fromAddr.Hash())
	if err != nil {
		fmt.Println("cannot resolve sender address")
		os.Exit(1)
	}
	recipientInfo, err := svc.ResolveAddress(toAddr.Hash())
	if err != nil {
		fmt.Println("cannot resolve recipient address")
		os.Exit(1)
	}

	// Setup addressing
	addressing := message.NewAddressing(
		*fromAddr,
		privKey,
		senderInfo.RoutingInfo.Routing,
		*toAddr,
		&recipientInfo.PublicKey,
	)

	// Compose mail
	envelope, err := message.Compose(addressing, opts.Subject, opts.Blocks, opts.Attachments)
	if err != nil {
		fmt.Println("cannot compose message")
		os.Exit(1)
	}

	// Send mail
	client, err := api.NewAuthenticated(*fromAddr, privKey, senderInfo.RoutingInfo.Routing, apiErrorFunc)
	if err != nil {
		fmt.Println("cannot connect to api")
		os.Exit(1)
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		fmt.Println("cannot send message: ", err)
		os.Exit(1)
	}
}

func parseFlags() {
	parser := flags.NewParser(&opts, flags.IgnoreUnknown)
	_, err := parser.Parse()

	if err != nil {
		flagsError, _ := err.(*flags.Error)
		if flagsError.Type == flags.ErrHelp {
			os.Exit(1)
		}

		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if opts.Message == "" && len(opts.Blocks) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Please specify either a message (-m) or one or more blocks (-b)")
		os.Exit(1)
	}

	if opts.Message != "" && len(opts.Blocks) > 0 {
		_, _ = fmt.Fprintln(os.Stderr, "Please specify either a message (-m) or one or more blocks (-b)")
		os.Exit(1)
	}

	// Check key
	privKey, err = bmcrypto.PrivateKeyFromString(opts.PrivateKey)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Incorrect private key specified.")
		os.Exit(1)
	}

	// Check from address
	fromAddr, err = address.NewAddress(opts.From)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Incorrect sender address specified.")
		os.Exit(1)
	}

	// Check to address
	toAddr, err = address.NewAddress(opts.To)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Incorrect recipient address specified.")
		os.Exit(1)
	}
}
