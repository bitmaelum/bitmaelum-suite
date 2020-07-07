package handlers

import (
	"fmt"
	"github.com/bitmaelum/bitmaelum-suite/core"
	"github.com/bitmaelum/bitmaelum-suite/core/api"
	"github.com/bitmaelum/bitmaelum-suite/core/container"
	"github.com/bitmaelum/bitmaelum-suite/core/encrypt"
	"github.com/bitmaelum/bitmaelum-suite/core/resolve"
	"github.com/bitmaelum/bitmaelum-suite/internal/account"
	"github.com/bitmaelum/bitmaelum-suite/internal/message"
	"github.com/bitmaelum/bitmaelum-suite/pkg/address"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"strings"
)

// ComposeMessage composes a new message from the given account Info to the "to" with given subject, blocks and attachments
func ComposeMessage(info account.Info, toAddr address.HashAddress, subject string, b, a []string) error {
	// Resolve public key for our recipient
	resolver := container.GetResolveService()
	toInfo, err := resolver.Resolve(toAddr)
	if err != nil {
		return fmt.Errorf("cannot retrieve public key for '%s'", toAddr.String())
	}

	// Generate blocks and attachments
	blocks, err := generateBlocks(b)
	if err != nil {
		return err
	}
	attachments, err := generateAttachments(a)
	if err != nil {
		return err
	}

	// Generate catalog
	catalog, err := generateCatalog(info, toAddr, subject, blocks, attachments)
	if err != nil {
		return err
	}

	// Encrypt catalog for upload
	catalogKey, encryptedCatalog, err := encrypt.CatalogEncrypt(*catalog)
	if err != nil {
		return err
	}

	// Generate header based on our encrypted catalog
	header, err := generateHeader(info, toInfo, encryptedCatalog, catalogKey)
	if err != nil {
		return err
	}

	msgID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	fmt.Printf("Outgoing message created: %s\n", msgID.String())
	fmt.Printf("Sending message to : %s\n", info.Server)

	err = uploadToServer(msgID.String(), info, header, encryptedCatalog, catalog)
	if err != nil {
		return err
	}
	return nil
}

func uploadToServer(msgID string, info account.Info, header *message.Header, encryptedCatalog []byte, catalog *message.Catalog) error {
	// Upload message to server
	addr, err := address.NewHash(info.Address)
	if err != nil {
		return err
	}

	client, err := api.New(&info)
	if err != nil {
		return err
	}

	// parallelize uploads
	g := new(errgroup.Group)
	g.Go(func() error {
		return client.UploadHeader(*addr, msgID, header)
	})
	g.Go(func() error {
		return client.UploadCatalog(*addr, msgID, encryptedCatalog)
	})
	for _, block := range catalog.Blocks {
		// Store locally, otherwise the anonymous go function doesn't know which "block"
		b := block
		g.Go(func() error {
			return client.UploadBlock(*addr, msgID, b.ID, b.Reader)
		})
	}
	for _, attachment := range catalog.Attachments {
		// Store locally, otherwise the anonymous go function doesn't know which "attachment"
		a := attachment
		g.Go(func() error {
			return client.UploadBlock(*addr, msgID, a.ID, a.Reader)
		})
	}

	// Wait until all are completed
	if err := g.Wait(); err != nil {
		_ = client.DeleteMessage(*addr, msgID)
		return err
	}

	// All done, mark upload as completed
	client.CompleteUpload(*addr, msgID)

	return nil
}

// Generate a header file based on the info provided
func generateHeader(info account.Info, toInfo *resolve.Info, catalog []byte, catalogKey []byte) (*message.Header, error) {
	header := &message.Header{}

	// We can add a multitude of checksums here.. whatever we like
	header.Catalog.Checksum = append(header.Catalog.Checksum, core.Sha256(catalog))
	header.Catalog.Checksum = append(header.Catalog.Checksum, core.Sha1(catalog))
	header.Catalog.Checksum = append(header.Catalog.Checksum, core.Crc32(catalog))
	header.Catalog.Checksum = append(header.Catalog.Checksum, core.Md5(catalog))
	header.Catalog.Size = uint64(len(catalog))
	header.Catalog.Crypto = "rsa+aes256gcm"

	pubKey, err := encrypt.PEMToPubKey([]byte(toInfo.PublicKey))
	if err != nil {
		return nil, err
	}

	header.Catalog.EncryptedKey, err = encrypt.Encrypt(pubKey, catalogKey)
	if err != nil {
		return nil, err
	}

	header.To.Addr = address.HashAddress(toInfo.Hash)

	h, err := address.NewHash(info.Address)
	if err != nil {
		return nil, err
	}
	header.From.Addr = *h

	header.From.PublicKey = info.PubKey
	header.From.ProofOfWork.Bits = info.Pow.Bits
	header.From.ProofOfWork.Proof = info.Pow.Proof

	return header, nil
}

// Generate a complete catalog file. Outputs catalog key and the encrypted catalog
func generateCatalog(info account.Info, toAddr address.HashAddress, subject string, b []message.Block, a []message.Attachment) (*message.Catalog, error) {
	// Create catalog
	cat := message.NewCatalog(&info)

	// @TODO: maybe these should be setters in Catalog?
	cat.To.Address = toAddr.String()
	cat.To.Name = toAddr.String()

	cat.Flags = append(cat.Flags, "important")
	cat.Labels = append(cat.Labels, "invoice", "sales", "seams-cms")
	cat.Subject = subject
	cat.ThreadID = ""

	for _, block := range b {
		err := cat.AddBlock(block)
		if err != nil {
			return nil, err
		}
	}
	for _, attachment := range a {
		err := cat.AddAttachment(attachment)
		if err != nil {
			return nil, err
		}
	}

	return cat, nil
}

// Generate message attachments based on the given paths to files
func generateAttachments(a []string) ([]message.Attachment, error) {
	// Parse attachments
	var attachments []message.Attachment
	for _, attachment := range a {
		_, err := os.Stat(attachment)
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("attachment %s does not exist", attachment)
		}

		reader, err := os.Open(attachment)
		if err != nil {
			return nil, fmt.Errorf("attachment %s cannot be opened", attachment)
		}

		attachments = append(attachments, message.Attachment{
			Path:   attachment,
			Reader: reader,
		})
	}

	return attachments, nil
}

// Generate message blocks based on the given strings
func generateBlocks(b []string) ([]message.Block, error) {
	// Parse blocks
	var blocks []message.Block
	for _, block := range b {
		split := strings.SplitN(block, ",", 2)
		if len(split) <= 1 {
			return nil, fmt.Errorf("please specify blocks in the format '<type>,<content>' or '<type>,file:<filename>'")
		}

		// By default assume content is inline
		size := int64(len(split[1]))
		var r io.Reader = strings.NewReader(split[1])

		if strings.HasPrefix(split[1], "file:") {
			// Open file as a reader
			f, err := os.Open(strings.TrimPrefix(split[1], "file:"))
			if err != nil {
				return nil, err
			}

			// Read file size
			fi, err := f.Stat()
			if err != nil {
				return nil, err
			}

			r = f
			size = fi.Size()
		}

		blocks = append(blocks, message.Block{
			Type:   split[0],
			Size:   uint64(size),
			Reader: r,
		})
	}

	return blocks, nil
}
