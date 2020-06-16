package main

import (
    "bytes"
    "fmt"
    "github.com/bitmaelum/bitmaelum-server/core"
    "github.com/bitmaelum/bitmaelum-server/core/account"
    "golang.org/x/crypto/ssh/terminal"
    "log"
    "syscall"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
    Addr        string      `short:"a" long:"address" description:"Address"`
}

var opts Options

func main() {
    core.ParseOptions(&opts)
    core.LoadClientConfig(opts.Config)

    var password []byte

    fmt.Print("\U0001F511 Enter current password (or enter when none)")
    password, _ = terminal.ReadPassword(int(syscall.Stdin))


    addr, err := core.NewAddressFromString(opts.Addr)
    if err != nil {
        log.Fatal(err)
    }
    ai, err := account.LoadAccount(*addr, password)
    if err != nil {
        log.Fatal(err)
    }

    for {
        fmt.Println("")
        fmt.Print("\U0001F511 Enter new password ")
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

    err = account.SaveAccount(*addr, password, *ai)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Account encrypted. Please do not loose your password!\n")
}
