package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/api"
	"github.com/bitmaelum/bitmaelum-server/core/checksum"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/core/encrypt"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"github.com/bitmaelum/bitmaelum-server/core/resolve"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func ComposeMessage(ai core.AccountInfo, to core.Address, subject string, b, a []string) error {
	// Resolve public key for our recipient
	resolver := container.GetResolveService()
	toInfo, err := resolver.Resolve(to)
	if err != nil {
		return fmt.Errorf("cannot retrieve public key for '%s'", to.String())
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
	catalog, err := generateCatalog(ai, toInfo, subject, blocks, attachments)
	if err != nil {
		return err
	}

	// Encrypt catalog for upload
	catalogKey, encryptedCatalog, err := encrypt.EncryptCatalog(*catalog)
	if err != nil {
		return err
	}

	// Generate header based on our encrypted catalog
	header, err := generateHeader(ai, toInfo, encryptedCatalog, catalogKey)
	if err != nil {
		return err
	}

	msgUuid, err := uploadToServer(ai, header, encryptedCatalog, catalog)
	if err != nil {
		return err
	}
	fmt.Printf("Message uploaded in %s\n", msgUuid)

	err = writeMessageToDisk(msgUuid, header, encryptedCatalog)
	if err != nil {
		return err
	}
	fmt.Printf("Message stored in %s\n", msgUuid)

	return nil
}

func uploadToServer(ai core.AccountInfo, header *message.Header, encryptedCatalog []byte, catalog *message.Catalog) (string, error) {
	// Upload message to server
	messageId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	addr := core.StringToHash(ai.Address)

	client, err := api.CreateNewClient(&ai)
	if err != nil {
		return "", err
	}

	// Upload header and catalog
	err = client.UploadHeader(addr, messageId.String(), header)
	if err != nil {
		_ = client.DeleteMessage(addr, messageId.String())
		return "", err
	}

	err = client.UploadCatalog(addr, messageId.String(), encryptedCatalog)
	if err != nil {
		_ = client.DeleteMessage(addr, messageId.String())
		return "", err
	}

	// Upload blocks & attachments
	for _, block := range catalog.Blocks {
		err = client.UploadBlock(addr, messageId.String(), block.Id, block.Reader)
		if err != nil {
			_ = client.DeleteMessage(addr, messageId.String())
			return "", err
		}
	}
	for _, attachment := range catalog.Attachments {
		err = client.UploadBlock(addr, messageId.String(), attachment.Id, attachment.Reader)
		if err != nil {
			_ = client.DeleteMessage(addr, messageId.String())
			return "", err
		}
	}

	return messageId.String(), nil
}

// Write the given message to disk
func writeMessageToDisk(msgUuid string, header interface{}, catalog []byte) error {
	p := ".out/" + msgUuid

	err := os.MkdirAll(p, 0755)
	if err != nil {
		return err
	}

	// Write catalog
	err = ioutil.WriteFile(p+"/catalog.json.enc", catalog, 0600)
	if err != nil {
		_ = os.RemoveAll(p)
		return err
	}

	// Write header
	data, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		_ = os.RemoveAll(p)
		return fmt.Errorf("error trying to marshal header: %v", err)
	}

	err = ioutil.WriteFile(".out/"+msgUuid+"/header.json", data, 0600)
	if err != nil {
		_ = os.RemoveAll(p)
		return err
	}

	return nil
}

// Generate a header file based on the info provided
func generateHeader(ai core.AccountInfo, to *resolve.ResolveInfo, catalog []byte, catalogKey []byte) (*message.Header, error) {
	header := &message.Header{}

	// We can add a multitude of checksums here.. whatever we like
	header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Sha256(catalog))
	header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Sha1(catalog))
	header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Crc32(catalog))
	header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Md5(catalog))
	header.Catalog.Size = uint64(len(catalog))
	header.Catalog.Crypto = "rsa+aes256gcm"

	pubKey, err := encrypt.PEMToPubKey([]byte(to.PublicKey))
	if err != nil {
		return nil, err
	}

	header.Catalog.EncryptedKey, err = encrypt.Encrypt(pubKey, catalogKey)
	if err != nil {
		return nil, err
	}

	header.To.Addr = core.StringToHash(to.Address)

	header.From.Addr = core.StringToHash(ai.Address)
	header.From.PublicKey = ai.PubKey
	header.From.ProofOfWork.Bits = ai.Pow.Bits
	header.From.ProofOfWork.Proof = ai.Pow.Proof

	return header, nil
}

// Generate a complete catalog file. Outputs catalog key and the encrypted catalog
func generateCatalog(ai core.AccountInfo, to *resolve.ResolveInfo, subject string, b []message.Block, a []message.Attachment) (*message.Catalog, error) {
	// Create catalog
	cat := message.NewCatalog(&ai)

	// @TODO: maybe these should be setters in Catalog?
	cat.To.Address = to.Address
	cat.To.Name = ""

	cat.Flags = append(cat.Flags, "important")
	cat.Labels = append(cat.Labels, "invoice", "sales", "seams-cms")
	cat.Subject = subject
	cat.ThreadId = ""

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
		split := strings.Split(block, ",")
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
