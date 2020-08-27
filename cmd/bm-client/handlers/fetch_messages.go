package handlers

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/pkg"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/c2h5oh/datasize"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

// FetchMessages will display message information from a box or display all boxes
func FetchMessages(info *pkg.Info, box string, checkOnly bool) {
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

	if box == "" || box == "0" {
		displayBoxList(client, *addr)
	} else {
		displayBox(client, *addr, info, box)
	}
}

func displayBoxList(client *api.API, addr address.HashAddress) {
	mbl, err := client.GetMailboxList(addr)
	if err != nil {
		logrus.Fatal(err)
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

func displayBox(client *api.API, addr address.HashAddress, info *pkg.Info, box string) {
	mb, err := client.GetMailboxMessages(addr, box)
	if err != nil {
		logrus.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"ID", "Subject", "From", "Organisation", "Date", "# Blocks", "# Attachments"}
	table.SetHeader(headers)

	for _, msg := range mb.Messages {
		key, err := encrypt.Decrypt(info.PrivKey, msg.Header.Catalog.EncryptedKey)
		if err != nil {
			logrus.Fatal(err)
		}
		catalog, err := encrypt.CatalogDecrypt(key, msg.Catalog)
		if err != nil {
			continue
		}

		blocks := []string{}
		for _, b := range catalog.Blocks {
			blocks = append(blocks, b.Type)
		}

		attachments := []string{}
		for _, a := range catalog.Attachments {
			fs := datasize.ByteSize(a.Size)
			attachments = append(attachments, a.FileName+" ("+fs.HR()+")")
		}

		values := []string{
			msg.ID,
			catalog.Subject,
			catalog.From.Name,
			catalog.From.Organisation,
			catalog.CreatedAt.Format(time.RFC822),
			strings.Join(blocks, ","),
			strings.Join(attachments, "\n"),
			catalog.To.Name,
		}

		table.Append(values)
	}
	table.Render()
}
