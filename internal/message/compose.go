package message

import (
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
)

func Compose(envConfig Addressing, subject string, b, a []string) (*Envelope, error) {
	cat, err := generateCatalog(envConfig.Sender.Address, envConfig.Recipient.Address, subject, b, a)
	if err != nil {
		return nil, err
	}

	header, err := generateHeader(envConfig.Sender.Address, envConfig.Recipient.Address)
	if err != nil {
		return nil, err
	}

	envelope, err := header.NewEnvelope()
	if err != nil {
		return nil, err
	}

	err = envelope.AddCatalog(cat)
	if err != nil {
		return nil, err
	}

	err = envelope.AddHeader(header)
	if err != nil {
		return nil, err
	}

	// Close the envelope for sending
	err = envelope.CloseAndEncrypt(envConfig.Sender.PrivKey, envConfig.Recipient.PubKey)
	if err != nil {
		return nil, err
	}

	return envelope, nil
}


// Generate a header file based on the info provided
func generateHeader(sender, recipient address.Address) (*Header, error) {
	header := &Header{}

	header.To.Addr = sender.Hash()
	header.From.Addr = recipient.Hash()

	return header, nil
}

// Generate a catalog filled with blocks and attachments
func generateCatalog(sender, recipient address.Address, subject string, b, a []string) (*Catalog, error) {
	// Create a new catalog
	cat := NewCatalog(&sender, &recipient, subject)

	// Add blocks to catalog
	blocks, err := GenerateBlocks(b)
	if err != nil {
		return nil, err
	}
	for _, block := range blocks {
		err := cat.AddBlock(block)
		if err != nil {
			return nil, err
		}
	}

	// Add attachments to catalog
	attachments, err := GenerateAttachments(a)
	if err != nil {
		return nil, err
	}
	for _, attachment := range attachments {
		err := cat.AddAttachment(attachment)
		if err != nil {
			return nil, err
		}
	}

	return cat, nil
}
