package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/config"
	"github.com/bitmaelum/bitmaelum-suite/internal/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/resolver"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
	"github.com/bitmaelum/bitmaelum-suite/pkg/hash"
	"github.com/c2h5oh/datasize"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

// FetchMessages will display message information from a box or display all boxes
func FetchMessages(info *internal.AccountInfo, routingInfo *resolver.RoutingInfo, box string) {
	client, err := api.NewAuthenticated(info, api.ClientOpts{
		Host:          routingInfo.Routing,
		AllowInsecure: config.Client.Server.AllowInsecure,
		Debug:         config.Client.Server.DebugHTTP,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	addr := hash.New(info.Address)
	if box == "" || box == "0" {
		displayBoxList(client, addr)
	} else {
		displayBox(client, addr, info, box)
	}
}

func displayBoxList(client *api.API, addr hash.Hash) {
	mbl, err := client.GetMailboxList(addr)
	if err != nil {
		logrus.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"Mailbox ID", "Total Messages"}
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

func displayBox(client *api.API, addr hash.Hash, info *internal.AccountInfo, box string) {
	mb, err := client.GetMailboxMessages(addr, box)
	if err != nil {
		logrus.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	headers := []string{"ID", "Subject", "From", "Date", "# Blocks", "# Attachments"}
	table.SetHeader(headers)

	for _, msg := range mb.Messages {
		key, err := bmcrypto.Decrypt(info.PrivKey, msg.Header.Catalog.EncryptedKey)
		if err != nil {
			logrus.Fatal(err)
		}
		catalog := &message.Catalog{}
		err = encrypt.CatalogDecrypt(key, msg.Catalog, catalog)
		if err != nil {
			continue
		}

		var blocks []string
		for _, b := range catalog.Blocks {
			blocks = append(blocks, b.Type)
		}

		var attachments []string
		for _, a := range catalog.Attachments {
			fs := datasize.ByteSize(a.Size)
			attachments = append(attachments, a.FileName+" ("+fs.HR()+")")
		}

		values := []string{
			msg.ID,
			catalog.Subject,
			fmt.Sprintf("%s <%s>", catalog.From.Name, catalog.From.Address),
			catalog.CreatedAt.Format(time.RFC822),
			strings.Join(blocks, ","),
			strings.Join(attachments, "\n"),
		}

		table.Append(values)
	}
	table.Render()
}
