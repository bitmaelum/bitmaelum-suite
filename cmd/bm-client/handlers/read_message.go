package handlers

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/sirupsen/logrus"
)

// ReadMessage will read a specific message blocks
func ReadMessage(info *pkg.Info, box, messageID, block string) {
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
	message, err := client.GetMessage(*addr, box, messageID)
	if err != nil {
		logrus.Fatal(err)
	}

	// Get private key, and decrypt the encryption-key for the catalog
	privKey, err := encrypt.PEMToPrivKey([]byte(info.PrivKey))
	if err != nil {
		logrus.Fatal(err)
	}

	key, err := encrypt.Decrypt(privKey, message.Header.Catalog.EncryptedKey)
	if err != nil {
		logrus.Fatal(err)
	}

	// Decrypt the catalog
	catalog, err := encrypt.CatalogDecrypt(key, message.Catalog)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("%#v", catalog)
}
