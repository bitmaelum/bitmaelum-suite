package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
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
	addr, err := address.NewAddress(info.Address)
	if err != nil {
		logrus.Fatal(err)
	}
	msg, err := client.GetMessage(addr.Hash(), box, messageID)
	if err != nil {
		logrus.Fatal(err)
	}

	key, err := bmcrypto.Decrypt(info.PrivKey, msg.Header.Catalog.Crypto, msg.Header.Catalog.EncryptedKey)
	if err != nil {
		logrus.Fatal(err)
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

		data, err := client.GetMessageBlock(addr.Hash(), box, messageID, b.ID)
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
