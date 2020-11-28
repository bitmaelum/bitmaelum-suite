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
	"io/ioutil"
	"os"

	"github.com/bitmaelum/bitmaelum-suite/cmd/bm-client/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/internal/vault"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/c2h5oh/datasize"
	"github.com/sirupsen/logrus"
)

// ReadMessage will read a specific message blocks
func ReadMessage(info *vault.AccountInfo, routingInfo *resolver.RoutingInfo, box, messageID string, saveAttachments bool) {
	client, err := api.NewAuthenticated(*info.Address, &info.PrivKey, routingInfo.Routing, internal.JwtErrorFunc)
	if err != nil {
		logrus.Fatal(err)
	}

	// Fetch message from API
	msg, err := client.GetMessage(info.Address.Hash(), box, messageID)
	if err != nil {
		logrus.Fatal(err)
	}

	key, err := bmcrypto.Decrypt(info.PrivKey, msg.Header.Catalog.TransactionID, msg.Header.Catalog.EncryptedKey)
	if err != nil {
		logrus.Fatal(err)
	}

	// Verify the clientSignature
	if !message.VerifyClientHeader(msg.Header) {
		logrus.Fatalf("message %s has failed the client signature check. Seems that this message may have been spoofed.", messageID)
	}

	// Decrypt the catalog
	catalog := &message.Catalog{}
	err = bmcrypto.CatalogDecrypt(key, msg.Catalog, catalog)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("--------------------------------------------------------\n")
	fmt.Printf("From       : %s <%s>\n", catalog.From.Name, catalog.From.Address)
	fmt.Printf("To         : %s <%s>\n", catalog.To.Name, catalog.To.Address)
	fmt.Printf("Subject    : %s\n", catalog.Subject)
	fmt.Printf("\n")
	fmt.Printf("Msg ID     : %s\n", messageID)
	fmt.Printf("Created at : %s\n", catalog.CreatedAt)
	fmt.Printf("ThreadID   : %s\n", catalog.ThreadID)
	fmt.Printf("Flags      : %s\n", catalog.Flags)
	fmt.Printf("Labels     : %s\n", catalog.Labels)
	fmt.Println("--------------------------------------------------------")
	for idx, b := range catalog.Blocks {
		fmt.Printf("Block %02d: %-20s %8s\n", idx, b.Type, datasize.ByteSize(b.Size))
		fmt.Printf("\n")

		data, err := client.GetMessageBlock(info.Address.Hash(), box, messageID, b.ID)
		if err != nil {
			continue
		}
		bb := bytes.NewBuffer(data)

		r, err := bmcrypto.GetAesDecryptorReader(b.IV, b.Key, bb)
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Println(string(content))
	}
	fmt.Println("--------------------------------------------------------")
	for idx, b := range catalog.Attachments {
		fmt.Printf("Attachment %02d: %30s %8d %s\n", idx, b.FileName, datasize.ByteSize(b.Size), b.MimeType)

		if saveAttachments {
			ar, err := client.GetMessageAttachment(info.Address.Hash(), box, messageID, b.ID)
			if err == nil {
				err = saveAttachment(b, ar)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	fmt.Println("--------------------------------------------------------")
}

// saveAttachment will save the given attachment to disk
func saveAttachment(att message.AttachmentType, ar io.ReadCloser) error {
	defer func() {
		_ = ar.Close()
	}()

	_, ok := os.Stat(att.FileName)
	if ok == nil {
		fmt.Printf("cannot write to %s: file exists\n", att.FileName)
		return errors.New("cannot write to file")
	}

	r, err := bmcrypto.GetAesDecryptorReader(att.IV, att.Key, ar)
	if err != nil {
		fmt.Printf("cannot create decryptor to %s: %s\n", att.FileName, err)
		return err
	}

	f, err := os.Create(att.FileName)
	if err != nil {
		fmt.Printf("cannot open file %s: %s\n", att.FileName, err)
		return err
	}

	if att.Compression == "zlib" {
		r, err = message.ZlibDecompress(r)
		if err != nil {
			fmt.Printf("error while creating zlib reader %s: %s\n", att.FileName, err)
			return err
		}
	}

	n, err := io.Copy(f, r)
	if err != nil || n != int64(att.Size) {
		fmt.Printf("error while writing file %s: %s (%d/%d bytes)\n", att.FileName, err, n, att.Size)
		return err
	}

	_ = f.Close()
	fmt.Printf("saved file: %s\n", att.FileName)

	return nil

}
