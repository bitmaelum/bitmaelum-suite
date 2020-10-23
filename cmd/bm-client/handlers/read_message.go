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
	"fmt"
	"io/ioutil"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	clientSig "github.com/bitmaelum/bitmaelum-suite/internal/client"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/c2h5oh/datasize"
	"github.com/sirupsen/logrus"
)

// ReadMessage will read a specific message blocks
func ReadMessage(info *internal.AccountInfo, routingInfo *resolver.RoutingInfo, box, messageID, blockType string) {
	client, err := api.NewAuthenticated(info, api.ClientOpts{
		Host:          routingInfo.Routing,
		AllowInsecure: config.Client.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// Fetch message from API
	msg, err := client.GetMessage(info.AddressHash(), box, messageID)
	if err != nil {
		logrus.Fatal(err)
	}

	key, err := bmcrypto.Decrypt(info.PrivKey, msg.Header.Catalog.TransactionID, msg.Header.Catalog.EncryptedKey)
	if err != nil {
		logrus.Fatal(err)
	}

	// Verify the clientSignature
	if !clientSig.VerifyHeader(msg.Header) {
		logrus.Fatalf("message %s has failed the client signature check. Seems that this message may have been spoofed.", messageID)
	}

	// Decrypt the catalog
	catalog := &message.Catalog{}
	err = encrypt.CatalogDecrypt(key, msg.Catalog, catalog)
	if err != nil {
		logrus.Fatal(err)
	}

	// spew.Dump(catalog)
	//
	// for _, b := range catalog.Blocks {
	// 	data, err := client.GetMessageBlock(*addr, box, messageID, b.ID)
	// 	bb := bytes.NewBuffer(data)
	//
	// 	r, err := encrypt.GetAesDecryptorReader(b.IV, b.Key, bb)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	content, err := ioutil.ReadAll(r)
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	//
	// 	spew.Dump(b)
	// 	spew.Dump(content)
	// }
	//
	// for _, a := range catalog.Attachments {
	// 	ar, err := client.GetMessageAttachment(*addr, box, messageID, a.ID)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	r, err := encrypt.GetAesDecryptorReader(a.IV, a.Key, ar)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	content, err := ioutil.ReadAll(r)
	// 	if err != nil {
	// 		logrus.Fatal(err)
	// 	}
	//
	// 	spew.Dump(a)
	// 	spew.Dump(content)
	// }

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

		data, err := client.GetMessageBlock(info.AddressHash(), box, messageID, b.ID)
		if err != nil {
			panic(err)
		}
		bb := bytes.NewBuffer(data)

		r, err := encrypt.GetAesDecryptorReader(b.IV, b.Key, bb)
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Print(string(content))
	}
	fmt.Println("--------------------------------------------------------")
	for idx, b := range catalog.Attachments {
		fmt.Printf("Attachment %02d: %30s %8d %s\n", idx, b.FileName, datasize.ByteSize(b.Size), b.MimeType)
	}
	fmt.Println("--------------------------------------------------------")
}
