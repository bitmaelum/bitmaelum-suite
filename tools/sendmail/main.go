package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/core"
	"github.com/bitmaelum/bitmaelum-server/core/account/server"
	"github.com/bitmaelum/bitmaelum-server/core/checksum"
	"github.com/bitmaelum/bitmaelum-server/core/container"
	"github.com/bitmaelum/bitmaelum-server/core/encrypt"
	"github.com/bitmaelum/bitmaelum-server/core/message"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
    From        string      `short:"f" long:"from" description:"Sender address"`
    To          string      `short:"t" long:"to" description:"Recipient address"`
    Subject     string      `short:"s" long:"subject" description:"Subject of the message"`
    Block       []string    `short:"b" long:"block" description:"Content block"`
    Attachment  []string    `short:"a" long:"attachment" description:"Attachment"`
    Password    string      `short:"p" long:"password" description:"Password to decrypt your account"`
}

var opts Options

//var alreadyStoredInKeychain = false

func main() {
    core.ParseOptions(&opts)
    core.LoadClientConfig(opts.Config)


    // Convert strings into addresses
    fromAddr, err := core.NewAddressFromString(opts.From)
    if err != nil {
        log.Fatal(err)
    }
    toAddr, err := core.NewAddressFromString(opts.To)
    if err != nil {
        log.Fatal(err)
    }


    var pwd = []byte(opts.Password)

    //// @TODO: We don't want specific stuff here. We should have a more generic function to fetch passwords from either
    //// the OS's keychain (or equivalent), from the commandline arguments, or through console stdio
    //if runtime.GOOS == `darwin` {
    //    // Fetch passwd from keychain if not
    //    if len(opts.Password) == 0 {
    //        kc := &keychain.OSXKeyChain{}
    //        pwd, err = kc.Fetch(*fromAddr)
    //        if err == nil {
    //            alreadyStoredInKeychain = true
    //        }
    //    }
    //}

    // Load our FROM account
    ai, err := server.LoadAccount(*fromAddr, pwd)
    if err != nil {
        log.Fatal(err)
    }

    //// @TODO: We don't want specific stuff here. We should have a more generic function to fetch passwords from either
    //// the OS's keychain (or equivalent), from the commandline arguments, or through console stdio
    //// Store password in keychain
    //if runtime.GOOS == `darwin` {
    //    if !alreadyStoredInKeychain {
    //        kc := &keychain.OSXKeyChain{}
    //        err = kc.Store(*fromAddr, pwd)
    //        if err != nil {
    //            log.Fatal(err)
    //        }
    //    }
    //}


    // Resolve public key for our recipient
    resolver := container.GetResolveService()
    resolvedInfo, err := resolver.Resolve(*toAddr)
    if err != nil {
        log.Fatal(fmt.Sprintf("cannot retrieve public key for '%s'", opts.To))
    }


    // Create message id and temporary outbox
    msgUuid, err := uuid.NewRandom()
    if err != nil {
        log.Fatal(err)
    }
    err = os.MkdirAll(".out/" + msgUuid.String(), 0755)
    if err != nil {
        log.Fatal(err)
    }

    // Parse blocks
    var blocks []message.Block
    for idx := range opts.Block {
        split := strings.Split(opts.Block[idx], ",")
        if len(split) <= 1 {
            log.Fatal("Please specify blocks in the format '<type>,<content>' or '<type>,file:<filename>'")
        }


        // By default assume content is inline
        size := int64(len(split[1]))
        var r io.Reader = strings.NewReader(split[1])

        if (strings.HasPrefix(split[1], "file:")) {
            // Open file as a reader
            f, err := os.Open(strings.TrimPrefix(split[1], "file:"))
            if err != nil {
                log.Fatal(err)
            }

            // Read file size
            fi, err := f.Stat()
            if err != nil {
                log.Fatal(err)
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



    // Parse attachments
    var attachments []message.Attachment
    for idx := range opts.Attachment {
        _, err := os.Stat(opts.Attachment[idx])
        if os.IsNotExist(err) {
            log.Fatal(fmt.Sprintf("attachment %s does not exist", opts.Attachment[idx]))
        }

        reader, err := os.Open(opts.Attachment[idx])
        if err != nil {
            log.Fatal(fmt.Sprintf("attachment %s or cannot be opened", opts.Attachment[idx]))
        }

        attachments = append(attachments, message.Attachment{
            Path: opts.Attachment[idx],
            Reader: reader,
        })
    }



    // Create catalog
    cat := message.NewCatalog(ai)

    cat.To.Address = opts.To
    cat.To.Name = ""

    cat.Flags = append(cat.Flags, "important")
    cat.Labels = append(cat.Labels, "invoice", "sales", "seams-cms")
    cat.Subject = opts.Subject
    cat.ThreadId = ""


    for idx := range blocks {
       err = cat.AddBlock(blocks[idx])
       if err != nil {
           log.Fatal(err)
       }
    }
    for idx := range attachments {
       err = cat.AddAttachment(attachments[idx])
       if err != nil {
           log.Fatal(err)
       }
    }

    data, _ := json.MarshalIndent(cat, "", "  ")
    _ = ioutil.WriteFile(".out/" + msgUuid.String() + "/catalog.json", data, 0600)

    catalogKey, encCatalog, err := encrypt.EncryptCatalog(*cat)
    if err != nil {
        log.Fatalf("Error while encrypting catalog: %v", err)
    }

    _ = ioutil.WriteFile(".out/" + msgUuid.String() + "/catalog.json.enc", encCatalog, 0600)



    header := CreateHeader()

    // We can add a multitude of checksums here.. whatever we like
    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Sha256(encCatalog))
    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Sha1(encCatalog))
    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Crc32(encCatalog))
    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Md5(encCatalog))
    header.Catalog.Size = uint64(len(encCatalog))
    header.Catalog.Crypto = "rsa+aes256gcm"
    header.Catalog.EncryptedKey, err = encrypt.Encrypt([]byte(resolvedInfo.PublicKey), catalogKey)
    if err != nil {
        log.Fatalf("trying to encrypt keys: %v", err)
    }

    header.To.Addr = toAddr.Hash()

    header.From.Addr = fromAddr.Hash()
    header.From.PublicKey = ai.PubKey
    header.From.ProofOfWork.Bits = ai.Pow.Bits
    header.From.ProofOfWork.Proof = ai.Pow.Proof

    data, err = json.MarshalIndent(header, "", "  ")
    if err != nil {
        log.Fatalf("error trying to marshal header: %v", err)
    }

    _ = ioutil.WriteFile(".out/" + msgUuid.String() + "/header.json", data, 0600)
}


func CreateHeader() *message.Header {
    hdr := &message.Header{}
    return hdr;
}
