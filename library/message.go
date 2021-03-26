package bitmaelumClient

import (
	"github.com/bitmaelum/bitmaelum-suite/internal/api"
	"github.com/bitmaelum/bitmaelum-suite/internal/container"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/internal/messages"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/bitmaelum/bitmaelum-suite/pkg/bmcrypto"
)

func (b *BitMaelumClient) SendSimpleMessage(fromAcc, fromName, to, privKey, subject, body string) error {
	svc := container.Instance.GetResolveService()

	// Check sender
	fromAddr, err := address.NewAddress(fromAcc)
	if err != nil {
		return err
	}

	senderInfo, _ := svc.ResolveAddress(fromAddr.Hash())

	// Check recipient
	toAddr, err := address.NewAddress(to)
	if err != nil {
		return err
	}

	recipientInfo, err := svc.ResolveAddress(toAddr.Hash())
	if err != nil {
		return err
	}

	// Convert privKey string to bmcrypto
	pk, err := bmcrypto.PrivateKeyFromString(privKey)
	if err != nil {
		return err
	}

	// Setup addressing
	addressing := message.NewAddressing(message.SignedByTypeOrigin)
	addressing.AddSender(fromAddr, nil, fromName, *pk, senderInfo.RoutingInfo.Routing)
	addressing.AddRecipient(toAddr, nil, &recipientInfo.PublicKey)

	// Setup blocks
	var blocks []string
	blocks = append(blocks, "default,"+body)

	// Compose mail
	envelope, err := message.Compose(addressing, subject, blocks, nil)
	if err != nil {
		return err
	}

	// Send mail
	client, err := api.NewAuthenticated(*fromAddr, *pk, senderInfo.RoutingInfo.Routing, nil)
	if err != nil {
		return err
	}

	err = messages.Send(*client, envelope)
	if err != nil {
		return err
	}

	return nil
}
