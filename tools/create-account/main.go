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
    "io"
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

    fmt.Println(core.GetAsciiLogo())
    fmt.Println()
    fmt.Println("Account generation for BitMaelum")
    fmt.Println("--------------------------------")
    fmt.Println("")

    fmt.Println("Please specify the account you want to create.")
    fmt.Println("Use the 'name@organisation!' or 'name!' format.")
    fmt.Println("")

    address := readFromInput(InputReaderOpts{
        r: os.Stdin,
        prefix: "\U0001F4E8 Enter address: ",
        error: "Please specify an address in the correct format",
        validate: func(s string) bool {
            return core.IsValidAddress(s)
        },
    })


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
    fmt.Println("This seems a valid and free address. Let's get some more information for your new account")

    fmt.Println("")
    name := readFromInput(InputReaderOpts{
        r: os.Stdin,
        prefix: "\U0001F464 Enter your full name: ",
    })

    fmt.Println("")
    mailserver := readFromInput(InputReaderOpts{
        r: os.Stdin,
        prefix: "\U0001F4BB What is the BitMaelum server you want to store your account on: ",
    })

    token := readFromInput(InputReaderOpts{
        r: os.Stdin,
        prefix: "\U0001F4BB Please enter your invitation token that you have received from the mailserver: ",
    })

    fmt.Println("")
    fmt.Println("\U0001F510 Let's generate a key-pair for our new account... (this might take a while)")
    privateKey, publicKey, err := encrypt.GenerateKeyPair(encrypt.KeyTypeRSA)
    if err != nil {
        panic(err)
    }


    fmt.Println("")
    fmt.Println("\U0001F477 Let's do some proof-of-work... (this might take a while)")
    proof := core.NewProofOfWork(23, addr.Bytes(), 0)

    acc := core.AccountInfo{
        Address:      address,
        Name:         name,
        Organisation: "",
        PrivKey:      string(privateKey),
        PubKey:       string(publicKey),
        Pow:          *proof,
        Server:       mailserver,
    }


    fmt.Println("\U0001F510 Your account must be protected with a password. Without this password, you are not able to read or send mail, so don't loose it!")
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


    // All info is available.
    // Step 1: create a local account

    err = account.CreateLocalAccount(*addr, password, acc)
    if err != nil {
    err = account.CreateRemoteAccount(*addr, token, acc)
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

type InputReaderOpts struct {
    r           io.Reader
    prefix      string
    error       string
    validate    func(s string) bool
}

func readFromInput(opts InputReaderOpts) string {
    reader := bufio.NewReader(opts.r)

    var s string
    for {
        fmt.Print(opts.prefix)
        s, _ = reader.ReadString('\n')
        s = strings.Trim(s, "\n")

        if opts.validate == nil || opts.validate(s) {
            return s
        }

        fmt.Println(opts.error)
    }
}
