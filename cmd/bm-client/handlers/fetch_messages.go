package handlers

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"time"
)

func FetchMessages(info *account.Info, box string, checkOnly bool) {
	client, err := api.NewAuthenticated(info, api.ClientOpts{
		Host:          info.Server,
		AllowInsecure: config.Client.Server.AllowInsecure,
	})
	if err != nil {
		panic(err)
	}

	addr, err := address.NewHash(info.Address)
	if err != nil {
		panic(err)
	}

	if box == "" || box == "0" {
		displayBoxList(client, *addr)
	} else {
		displayBox(client, *addr, info, box)
	}
}

func displayBoxList(client *api.API, addr address.HashAddress) {
	mbl, err := client.GetMailboxList(addr)
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Mailbox ID", "Total Messages"}
	// headers := []string{"Subject", "From", "Organisation", "Date"}
	table.SetHeader(headers)

	for _, mb := range mbl.Boxes {
		values := []string{
			strconv.Itoa(mb.ID),
			strconv.Itoa(mb.Total),
		}

		table.Append(values)
	}
	table.Render()
}

func displayBox(client *api.API, addr address.HashAddress, info *account.Info, box string) {
	mb, err := client.GetMailboxMessages(addr, box)
	if err != nil {
		panic(err)
	}

	privKey, err := encrypt.PEMToPrivKey([]byte(info.PrivKey))
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Subject", "From", "Organisation", "Date", "# Blocks", "# Attachments"}
	table.SetHeader(headers)

	for _, msg := range mb.Messages {
		key, err := encrypt.Decrypt(privKey, msg.Header.Catalog.EncryptedKey)
		if err != nil {
			panic(err)
		}
		catalog, err := encrypt.CatalogDecrypt(key, msg.Catalog)
		if err != nil {
			continue
		}

		values := []string{
			catalog.Subject,
			catalog.From.Name,
			catalog.From.Organisation,
			catalog.CreatedAt.Format(time.RFC822),
			strconv.Itoa(len(catalog.Blocks)),
			strconv.Itoa(len(catalog.Attachments)),
			catalog.To.Name,
		}

		table.Append(values)
	}
	table.Render()
}
