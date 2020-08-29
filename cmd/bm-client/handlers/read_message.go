package handlers

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// ReadMessage will read a specific message blocks
func ReadMessage(info *pkg.Info, box, messageID, blockType string) {
	client, err := api.NewAuthenticated(info, api.ClientOpts{
		Host:          info.Server,
		AllowInsecure: config.Client.Server.AllowInsecure,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	addr, err := address.NewHash(info.Address)
	if err != nil {
		logrus.Fatal(err)
	}

	// Fetch message from API
	msg, err := client.GetMessage(*addr, box, messageID)
	if err != nil {
		logrus.Fatal(err)
	}

	key, err := encrypt.Decrypt(info.PrivKey, msg.Header.Catalog.EncryptedKey)
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

	for _, a := range catalog.Attachments {
		ar, err := client.GetMessageAttachment(*addr, box, messageID, a.ID)
		if err != nil {
			panic(err)
		}

		r, err := encrypt.GetAesDecryptorReader(a.IV, a.Key, ar)
		if err != nil {
			panic(err)
		}

		content, err := ioutil.ReadAll(r)
		if err != nil {
			logrus.Fatal(err)
		}

		spew.Dump(a)
		spew.Dump(content)
	}
}
