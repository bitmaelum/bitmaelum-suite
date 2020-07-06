package cmd

import (
	"errors"
	"fmt"
	"github.com/bitmaelum/bitmaelum-server/cmd/bm-client/handlers"
	"github.com/bitmaelum/bitmaelum-server/internal/config"
	"github.com/bitmaelum/bitmaelum-server/pkg/address"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
)

var composeCmd = &cobra.Command{
	Use:     "compose",
	Aliases: []string{"write", "send"},
	Short:   "Compose a new message",
	Long:    `This command will allow you to compose a new message and send it through your BitMaelum server`,
	Run: func(cmd *cobra.Command, args []string) {

		fromAddr, err := address.New(*from)
		if err != nil {
			panic(err)
		}
		fromInfo, err := Vault.GetAccountInfo(*fromAddr)
		if err != nil {
			panic(err)
		}

		toAddr, err := address.NewHash(*to)
		if err != nil {
			panic(err)
		}

		// If no blocks are specified, we assume reading a single block from stdin
		if len(*blocks) == 0 {
			var block string

			// Check if we have set $EDITOR so we can use this as our editor
			if hasEditorConfigured() {
				block, err = useRegularEditor()
			} else {
				// fall back to stdEditor
				block, err = useStdinEditor()
			}
			if err != nil {
				panic(err)
			}
			if len(block) == 0 {
				fmt.Println("Warning: empty message body")
			} else {
				*blocks = append(*blocks, "default,"+block)
			}
		}

		fmt.Printf("Endresult\n")
		fmt.Printf("  From:    %s (%s)\n", fromInfo.Name, fromInfo.Address)
		fmt.Printf("  To:      %s\n", *to)
		fmt.Printf("  Subject: %s\n", *subject)
		for i, block := range *blocks {
			fmt.Printf("Block  #%d %s\n", i, block)
		}
		for i, attachment := range *attachments {
			fmt.Printf("Att.   #%d %s\n", i, attachment)
		}

		err = handlers.ComposeMessage(*fromInfo, *toAddr, *subject, *blocks, *attachments)
		if err != nil {
			panic(err)
		}
	},
}

func findEditorPath() (string, error) {
	var editorPath = config.Client.Composer.Editor
	if editorPath != "" {
		return editorPath, nil
	}

	editorPath = os.Getenv("EDITOR")
	if editorPath != "" {
		return editorPath, nil
	}

	return "", errors.New("cannot find editor")
}

func hasEditorConfigured() bool {
	_, err := findEditorPath()
	return err == nil
}

func useRegularEditor() (string, error) {
	p, err := findEditorPath()
	if err != nil {
		return "", err
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "bm-")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	c := exec.Command(p, tmpFile.Name())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	err = c.Run()
	if err != nil {
		fmt.Printf("%#v", err)
		return "", err
	}

	data, err := ioutil.ReadFile(tmpFile.Name())
	return string(data), err
}

func useStdinEditor() (string, error) {
	fmt.Print("\U00002709 Enter your message and press CTRL-D when done.\n")

	data, err := ioutil.ReadAll(os.Stdin)
	fmt.Printf("%#v", err)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

var from, to, subject *string
var blocks, attachments *[]string

func init() {
	rootCmd.AddCommand(composeCmd)

	from = composeCmd.Flags().StringP("from", "f", "", "Sender address")
	to = composeCmd.Flags().StringP("to", "t", "", "Recipient address")
	subject = composeCmd.Flags().StringP("subject", "s", "", "Subject of the message")
	blocks = composeCmd.Flags().StringArrayP("blocks", "b", []string{}, "Message blocks")
	attachments = composeCmd.Flags().StringArrayP("attachment", "a", []string{}, "Attachments")

	_ = composeCmd.MarkFlagRequired("from")
	_ = composeCmd.MarkFlagRequired("to")
	_ = composeCmd.MarkFlagRequired("subject")
}
