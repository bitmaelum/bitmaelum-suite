package main

import (
    "bufio"
    "bytes"
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/account"
    "github.com/bitmaelum/bitmaelum-server/core/container"
    "github.com/bitmaelum/bitmaelum-server/core/encrypt"
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
    core.ParseOptions(&opts)
    core.LoadClientConfig(opts.Config)

    fmt.Println("Account generation for BitMaelum")
    fmt.Println("--------------------------------")
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
    _, err := ks.Resolve(*addr)
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
    fmt.Print("\U0001F4BB What is the bitmaelum server you want to store your account on: ")
    mailserver, _ := reader.ReadString('\n')
    mailserver = strings.TrimSuffix(mailserver, "\n")

    fmt.Println("")
    fmt.Println("\U0001F510 Let's generate a key-pair for our new account... (this might take a while)")
    privateKey, publicKey, err := encrypt.GenerateKeyPair(encrypt.KeyTypeRsa)
    if err != nil {
        panic(err)
    }


    fmt.Println("")
    fmt.Println("\U0001F477 Let's do some proof-of-work... (this might take a while)")
    proof := core.NewProofOfWork(22, addr.Bytes(), 0)

    acc := core.AccountInfo{
        Address:      address,
        Name:         name,
        Organisation: organisation,
        PrivKey:      string(privateKey),
        PubKey:       string(publicKey),
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
