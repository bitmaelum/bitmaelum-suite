package examples

import (
    "flag"
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

type GenerateFlags struct {
    email string
}

type EncryptFlags struct {
    from    string
    to      string
    input   string
}

type DecryptFlags struct {
    from    string
    to      string
    input   string
}


func main() {
    var cmdAccount = &cobra.Command{
        Use:    "account",
        Short:  "Account information",
        Long:   "Manage your account",
        Args:

    }

  var rootCmd = &cobra.Command{Use: "app"}
  rootCmd.AddCommand(cmdPrint, cmdEcho)
  cmdEcho.AddCommand(cmdTimes)
  rootCmd.Execute()

    generateFlags := GenerateFlags{}
    accountCmd := flag.NewFlagSet("account", flag.ExitOnError)
    generateCmd := accountCmd.Bool("generate", false, "Generate")
    generateFlags.email = accountCmd.String("email", "", "email address")

    encryptFlags := EncryptFlags{}
    encryptCmd := flag.NewFlagSet("encrypt", flag.ExitOnError)
    encryptFlags.from = encryptCmd.String("from", "", "email address")
    encryptFlags.to = encryptCmd.String("to", "", "email address")
    encryptFlags.input = encryptCmd.String("input", "", "file to encrypt")

    decryptFlags := DecryptFlags{}
    decryptCmd := flag.NewFlagSet("decrypt", flag.ExitOnError)
    decryptFlags.from = decryptCmd.String("from", "", "email address")
    decryptFlags.to = decryptCmd.String("to", "", "email address")
    decryptFlags.input = decryptCmd.String("input", "", "file to decrypt")

    flag.Parse()

    switch os.Args[1] {
    case "account":
        accountCmd.Parse(os.Args[2:])
    case "encrypt":
        encryptCmd.Parse(os.Args[2:])
    case "decrypt":
        decryptCmd.Parse(os.Args[2:])
    default:
        fmt.Printf("%q is not a valid command.\n", os.Args[1])
        os.Exit(2)
    }
}
