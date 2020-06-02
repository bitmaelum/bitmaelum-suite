package main

import (
    "encoding/json"
    "fmt"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/account"
    "github.com/jaytaph/mailv2/core/catalog"
    "github.com/jaytaph/mailv2/core/checksum"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jaytaph/mailv2/core/container"
    "github.com/jaytaph/mailv2/core/encode"
    "github.com/jaytaph/mailv2/core/encrypt"
    "github.com/jaytaph/mailv2/core/message"
    "github.com/jessevdk/go-flags"
    "io/ioutil"
    "os"
    "path"
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

func main() {
    // Parse config flags
    parser := flags.NewParser(&opts, flags.Default)
    if _, err := parser.Parse(); err != nil {
        flagsError, _ := err.(*flags.Error)
        if flagsError.Type == flags.ErrHelp {
            return
        }
        fmt.Println()
        parser.WriteHelp(os.Stdout)
        fmt.Println()
        return
    }

    // Load configuration
    err := config.Client.LoadConfig(path.Clean(opts.Config))
    if err != nil {
        panic(err)
    }


    // Convert strings into addresses
    fromAddr, err := core.NewAddressFromString(opts.From)
    if err != nil {
        panic(err)
    }
    toAddr, err := core.NewAddressFromString(opts.To)
    if err != nil {
        panic(err)
    }


    // Load our FROM account
    ai, err := account.LoadAccount(*fromAddr, []byte(opts.Password))
    if err != nil {
        panic(err)
    }

    // Resolve public key for our recipient
    resolver := container.GetResolveService()
    resolvedInfo, err := resolver.Resolve(*toAddr)
    if err != nil {
        panic(fmt.Sprintf("cannot retrieve public key for '%s'", opts.To))
    }

    fmt.Printf("Public found for reciever: %s", string(resolvedInfo.PublicKey))



    // Parse blocks
    var blocks []catalog.Block
    for idx := range opts.Block {
        split := strings.Split(opts.Block[idx], ",")
        if len(split) <= 1 {
            panic("Please specify blocks in the format '<type>,<content>' or '<type>,file:<filename>'")
        }

        var inlineContent = true

        var content = []byte(split[1])
        if (strings.HasPrefix(split[1], "file:")) {
            var err error
            content, err = ioutil.ReadFile(strings.TrimPrefix(split[1], "file:"))
            if err != nil {
                panic(fmt.Sprintf("Cannot read contents of file '%s'", strings.TrimPrefix(split[1], "file:")))
            }
            inlineContent = false
        }

        blocks = append(blocks, catalog.Block{
            Type: split[0],
            Inline: inlineContent,
            Content: content,
        })

    }



    // Parse attachments
    var attachments []catalog.Attachment
    for idx := range opts.Attachment {
        _, err := os.Stat(opts.Attachment[idx])
        if os.IsNotExist(err) {
            panic(fmt.Sprintf("attachment %s does not exist", opts.Attachment[idx]))
        }

        reader, err := os.Open(opts.Attachment[idx])
        if err != nil {
            panic(fmt.Sprintf("attachment %s or cannot be opened", opts.Attachment[idx]))
        }

        attachments = append(attachments, catalog.Attachment{
            Path: opts.Attachment[idx],
            Reader: reader,
        })
    }



    // Create catalog
    catalog := catalog.NewCatalog(ai)

    catalog.To.Address = opts.To
    catalog.To.Name = ""

    catalog.Flags = append(catalog.Flags, "important")
    catalog.Labels = append(catalog.Flags, "invoice", "sales", "seams-cms")
    catalog.Subject = opts.Subject
    catalog.ThreadId = ""


    for idx := range blocks {
       catalog.AddBlock(blocks[idx])
    }
    for idx := range attachments {
       _ = catalog.AddAttachment(attachments[idx])
    }

    data, _ := json.MarshalIndent(catalog, "", "  ")
    _ = ioutil.WriteFile("catalog.json", data, 0600)

    catalogKey, catalogIv, encCatalog, err := encrypt.EncryptCatalog(*catalog)
    if err != nil {
        panic(fmt.Sprintf("Error while encrypting catalog: %s", err))
    }

    _ = ioutil.WriteFile("catalog.json.enc", encCatalog, 0600)




    header := CreateHeader()

    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Sha256(encCatalog))
    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Sha1(encCatalog))
    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Crc32(encCatalog))
    header.Catalog.Checksum = append(header.Catalog.Checksum, checksum.Md5(encCatalog))
    header.Catalog.Size = uint64(len(encCatalog))
    header.Catalog.Crypto = "rsa+aes256"
    header.Catalog.Iv = encode.Encode(catalogIv)
    header.Catalog.Key, err = encrypt.EncryptKey([]byte(resolvedInfo.PublicKey), catalogKey)
    if err != nil {
        panic(fmt.Sprintf("trying to encrypt keys: %s", err))
    }

    header.To.Addr = toAddr.Hash()

    header.From.Addr = fromAddr.Hash()
    header.From.PublicKey = ai.PubKey
    header.From.ProofOfWork.Bits = ai.Pow.Bits
    header.From.ProofOfWork.Proof = ai.Pow.Proof

    data, err = json.Marshal(header)
    if err != nil {
        panic(fmt.Sprintf("error trying to marshal header: %s", err))
    }

    _ = ioutil.WriteFile("header.json", data, 0600)




    //data, _ := ioutil.ReadFile("catalog.json.enc")
    //data, err = encode.Decode(data)
    //
    //dc := message.Catalog{}
    //err = encrypt.DecryptJson(catalogKey, catalogIv, data, &dc)
    //
    //fmt.Printf("%#v", dc)


    //* fetch public key from receiver(s)
    //* header
    //* catalog


    //header, body := message.NewMessage(opts.From, opts.To, opts.Subject, blocks, attachments)
}


func CreateHeader() *message.Header {
    hdr := &message.Header{}
    return hdr;
}
