package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
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
	catalogKey, encryptedCatalog, err := generateCatalog(ai, toInfo, subject, blocks, attachments)
	if err != nil {
		return err
	}

	// Generate header
	header, err := generateHeader(ai, toInfo, encryptedCatalog, catalogKey)
	if err != nil {
		return err
	}

	msgUuid, err := writeMessageToDisk(header, encryptedCatalog)
	fmt.Printf("Message stored in %s", msgUuid)
	return err
}

func writeMessageToDisk(header interface{}, catalog []byte) (string, error) {
    // Create message id and create temporary outbox
    msgUuid, err := uuid.NewRandom()
    if err != nil {
        return "", err
    }
    err = os.MkdirAll(".out/" + msgUuid.String(), 0755)
    if err != nil {
        return "", err
    }

    // Write catalog
    err = ioutil.WriteFile(".out/" + msgUuid.String() + "/catalog.json.enc", catalog, 0600)
    if err != nil {
    	// @TODO: Remove out directory
		return "", err
	}

	// Write header
    data, err := json.MarshalIndent(header, "", "  ")
    if err != nil {
    	// @TODO: Remove out directory
        return "", fmt.Errorf("error trying to marshal header: %v", err)
    }

    err = ioutil.WriteFile(".out/" + msgUuid.String() + "/header.json", data, 0600)
    if err != nil {
    	// @TODO: Remove out directory
    	return "", err
	}

    return msgUuid.String(), nil
}

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

func generateCatalog(ai core.AccountInfo, to *resolve.ResolveInfo, subject string, b []message.Block, a []message.Attachment) ([]byte, []byte, error) {
    // Create catalog
    cat := message.NewCatalog(&ai)

    cat.To.Address = to.Address
    cat.To.Name = ""

    cat.Flags = append(cat.Flags, "important")
    cat.Labels = append(cat.Labels, "invoice", "sales", "seams-cms")
    cat.Subject = subject
    cat.ThreadId = ""

    for _, block := range b {
       	err := cat.AddBlock(block)
       	if err != nil {
			return nil, nil, err
       	}
    }
    for _, attachment := range a {
       err := cat.AddAttachment(attachment)
       if err != nil {
           return nil, nil, err
       }
    }

    catalogKey, encCatalog, err := encrypt.EncryptCatalog(*cat)
    if err != nil {
		return nil, nil, err
	}

	return catalogKey, encCatalog, err
}

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
            Path: attachment,
            Reader: reader,
        })
    }

    return attachments, nil
}

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
            Type: split[0],
            Size: uint64(size),
            Reader: r,
        })
    }

    return blocks, nil
}
