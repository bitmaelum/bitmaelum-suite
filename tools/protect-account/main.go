package main

import (
    "bytes"
    "fmt"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/account"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/jessevdk/go-flags"
    "golang.org/x/crypto/ssh/terminal"
    "os"
    "syscall"
)

type Options struct {
    Config      string      `short:"c" long:"config" description:"Configuration file" default:"./client-config.yml"`
    Addr        string      `short:"a" long:"address" description:"Address"`
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


    var password []byte

    fmt.Print("\U0001F511 Enter current password (or enter when none)")
    password, _ = terminal.ReadPassword(int(syscall.Stdin))


    addr, err := core.NewAddressFromString(opts.Addr)
    if err != nil {
        panic(err)
    }
    ai, err := account.LoadAccount(*addr, password)
    if err != nil {
        panic(err)
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
        panic(err)
    }

    fmt.Printf("Account encrypted. Please do not loose your password!\n")
}
