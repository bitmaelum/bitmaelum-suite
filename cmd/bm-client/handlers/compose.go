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

package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/client"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"golang.org/x/sync/errgroup"
)

// ComposeMessage composes a new message from the given account Info to the "to" with given subject, blocks and attachments
func ComposeMessage(info internal.AccountInfo, toAddr address.Address, subject string, b, a []string) error {
	// Resolve public key for our recipient
	resolver := container.GetResolveService()
	toInfo, err := resolver.ResolveAddress(toAddr.Hash())
	if err != nil {
		return fmt.Errorf("cannot retrieve public key for '%s'", toAddr.String())
	}

	// Generate blocks and attachments
	blocks, err := generateBlocks(b)
	if err != nil {
		return err
	}
	attachments, err := generateAttachments(a)
	if err != nil {
		return err
	}

	// Generate catalog
	catalog, err := generateCatalog(info, toAddr, subject, blocks, attachments)
	if err != nil {
		return err
	}

	// Encrypt catalog for upload
	catalogKey, encryptedCatalog, err := encrypt.CatalogEncrypt(*catalog)
	if err != nil {
		return err
	}

	// Generate header based on our encrypted catalog
	header, err := generateHeader(info, toInfo, encryptedCatalog, catalogKey)
	if err != nil {
		return err
	}

	//Sign the header
	err = client.SignHeader(header, info.PrivKey)
	if err != nil {
		return err
	}

	// Fetch routing info for the SENDER, as we send to the sender's mailserver (not the recipient)
	routingInfo, err := resolver.ResolveRouting(info.RoutingID)
	if err != nil {
		return err
	}

	// and finally send
	err = uploadToServer(info, *routingInfo, header, encryptedCatalog, catalog)
	if err != nil {
		return err
	}
	return nil
}

func uploadToServer(info internal.AccountInfo, routingInfo resolver.RoutingInfo, header *message.Header, encryptedCatalog []byte, catalog *message.Catalog) error {
	client, err := api.NewAuthenticated(&info, api.ClientOpts{
		Host:          routingInfo.Routing,
		AllowInsecure: config.Client.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
	})
	if err != nil {
		return err
	}

	// Get upload ticket
	t, err := client.GetTicket(header.From.Addr, header.To.Addr, "")
	if err != nil {
		return errors.New("cannot get ticket from server: " + err.Error())
	}
	if !t.Valid {
		return errors.New("invalid ticket returned by server")
	}

	// parallelize uploads
	g := new(errgroup.Group)
	g.Go(func() error {
		return client.UploadHeader(*t, header)
	})
	g.Go(func() error {
		return client.UploadCatalog(*t, encryptedCatalog)
	})
	for _, block := range catalog.Blocks {
		// Store locally, otherwise the anonymous go function doesn't know which "block"
		b := block
		g.Go(func() error {
			return client.UploadBlock(*t, b.ID, b.Reader)
		})
	}
	for _, attachment := range catalog.Attachments {
		// Store locally, otherwise the anonymous go function doesn't know which "attachment"
		a := attachment
		g.Go(func() error {
			return client.UploadBlock(*t, a.ID, a.Reader)
		})
	}

	// Wait until all are completed
	if err := g.Wait(); err != nil {
		_ = client.DeleteMessage(*t)
		return err
	}

	// All done, mark upload as completed
	return client.CompleteUpload(*t)
}

// Generate a header file based on the info provided
func generateHeader(info internal.AccountInfo, toInfo *resolver.AddressInfo, catalog []byte, catalogKey []byte) (*message.Header, error) {
	header := &message.Header{}

	// We can add a multitude of checksums here.. whatever we like
	r := bytes.NewBuffer(catalog)
	var err error
	header.Catalog.Checksum, err = message.CalculateChecksums(r)
	if err != nil {
		return nil, err
	}
	header.Catalog.Size = uint64(len(catalog))
	header.Catalog.EncryptedKey, header.Catalog.TransactionID, header.Catalog.Crypto, err = bmcrypto.Encrypt(toInfo.PublicKey, catalogKey)
	if err != nil {
		return nil, err
	}

	h, err := hash.NewFromHash(toInfo.Hash)
	if err != nil {
		return nil, err
	}
	header.To.Addr = *h
	header.From.Addr = info.AddressHash()
	header.From.PublicKey = &info.PubKey
	header.From.ProofOfWork = info.Pow

	return header, nil
}

// Generate a complete catalog file. Outputs catalog key and the encrypted catalog
func generateCatalog(info internal.AccountInfo, toAddr address.Address, subject string, b []message.Block, a []message.Attachment) (*message.Catalog, error) {
	// Create catalog
	cat := message.NewCatalog(&info)

	cat.AddFlags("new")
	cat.AddLabels("important")
	cat.SetToAddress(toAddr, "John Doe")

	cat.Subject = subject
	cat.ThreadID = ""

	for _, block := range b {
		err := cat.AddBlock(block)
		if err != nil {
			return nil, err
		}
	}
	for _, attachment := range a {
		err := cat.AddAttachment(attachment)
		if err != nil {
			return nil, err
		}
	}

	return cat, nil
}

// Generate message attachments based on the given paths to files
func generateAttachments(a []string) ([]message.Attachment, error) {
	// Parse attachments
	var attachments []message.Attachment
	for _, attachment := range a {
		_, err := os.Stat(attachment)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("attachment %s does not exist", attachment)
		}

		reader, err := os.Open(attachment)
		if err != nil {
			return nil, fmt.Errorf("attachment %s cannot be opened", attachment)
		}

		attachments = append(attachments, message.Attachment{
			Path:   attachment,
			Reader: reader,
		})
	}

	return attachments, nil
}

// Generate message blocks based on the given strings
func generateBlocks(b []string) ([]message.Block, error) {
	// Parse blocks
	var blocks []message.Block
	for _, block := range b {
		split := strings.SplitN(block, ",", 2)
		if len(split) <= 1 {
			return nil, fmt.Errorf("please specify blocks in the format '<type>,<content>' or '<type>,file:<filename>'")
		}

		// By default assume content is inline
		size := int64(len(split[1]))
		var r io.Reader = strings.NewReader(split[1])

		if strings.HasPrefix(split[1], "file:") {
			// Open file as a reader
			f, err := os.Open(strings.TrimPrefix(split[1], "file:"))
			if err != nil {
				return nil, err
			}

			// Read file size
			fi, err := f.Stat()
			if err != nil {
				return nil, err
			}

			r = f
			size = fi.Size()
		}

		blocks = append(blocks, message.Block{
			Type:   split[0],
			Size:   uint64(size),
			Reader: r,
		})
	}

	return blocks, nil
}
