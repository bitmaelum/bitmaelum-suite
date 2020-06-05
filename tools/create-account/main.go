package main

import (
    "bufio"
    "bytes"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/account"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jaytaph/mailv2/core/container"
    "github.com/jessevdk/go-flags"
    "golang.org/x/crypto/ssh/terminal"
    "os"
    "strings"
    "syscall"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
}

var opts Options

func main() {
    // Parse config
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
    err := config.Client.LoadConfig(opts.Config)
    if err != nil {
        panic(err)
    }

    fmt.Println("Account generation for mailv2")
    fmt.Println("-----------------------------")
    fmt.Println("")

    fmt.Println("Please specify the account you want to create.")
    fmt.Println("Use the 'name@organisation!' or 'name!' format.")
    fmt.Println("")

    reader := bufio.NewReader(os.Stdin)

    var address string
    for {
        fmt.Print("\U0001F4E8 Enter address: ")
        address, _ = reader.ReadString('\n')
        address = strings.Trim(address, "\n")

        if core.IsValidAddress(address) {
            break;
        }

        fmt.Println("Please specify an address in the correct format")
    }

    fmt.Println("")
    fmt.Printf("Checking if '%s' already exists...\n", address)

    addr, _ := core.NewAddressFromString(address)

    ks := container.GetResolveService()
    _, err = ks.Resolve(*addr)
    if err == nil {
        fmt.Printf("It seems that this address is already in use. Please specify another address.")
        os.Exit(1)
    }



    fmt.Println("")
    fmt.Println("This seems a valid and free address.")

    fmt.Println("")
    fmt.Println("Let's get some more information for your new account")

    fmt.Println("")
    fmt.Print("\U0001F464 Enter your full name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSuffix(name, "\n")
    fmt.Println("")

    fmt.Print("\U0001F3E2 Enter your organisation name, or leave empty when you have none: ")
    organisation, _ := reader.ReadString('\n')
    organisation = strings.TrimSuffix(organisation, "\n")

    fmt.Println("")
    fmt.Print("\U0001F4BB What is the mailv2 server you want to store your account on: ")
    mailserver, _ := reader.ReadString('\n')
    mailserver = strings.TrimSuffix(mailserver, "\n")

    fmt.Println("")
    fmt.Println("\U0001F510 Let's generate a key-pair for our new account... (this might take a while)")
    privateKey, err := core.CreateNewKeyPair(4096)
    if err != nil {
        panic(err)
    }

    pemData := x509.MarshalPKCS1PrivateKey(privateKey)
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: pemData,
    })

    pemData = x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: pemData,
    })

    fmt.Println("")
    fmt.Println("\U0001F477 Let's do some proof-of-work... (this might take a while)")
    proof := core.NewProofOfWork(22, addr.Bytes(), 0)

    acc := core.AccountInfo{
        Address:      address,
        Name:         name,
        Organisation: organisation,
        PrivKey:      string(privateKeyPEM),
        PubKey:       string(publicKeyPEM),
        Pow:          *proof,
    }


    var password []byte
    for {
        fmt.Println("")
        fmt.Print("\U0001F511 Enter your password ")
        password, _ = terminal.ReadPassword(int(syscall.Stdin))
        fmt.Println("")

        fmt.Println("")
        fmt.Print("\U0001F511 Retype your password ")
        password2, _ := terminal.ReadPassword(int(syscall.Stdin))
        fmt.Println("")

        if bytes.Equal(password, password2) {
            break
        }

        fmt.Println("passwords do not match. Please try again.")
    }

    err = account.SaveAccount(*addr, password, acc)
    if err != nil {
        panic(err)
    }

    fmt.Println("")
    fmt.Println("\U0001F310 Uploading resolve information and public key to the central resolve server")
    err = container.GetResolveService().UploadInfo(acc, mailserver)
    if err != nil {
        panic(err)
    }

    fmt.Println("")
    fmt.Println("-----------------------")
    fmt.Printf("Your new account configuration is stored. Make sure you do not forget your password!\n")
}
