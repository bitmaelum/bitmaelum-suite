package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AskMnemonicPhrase asks for the mnemonic phrase
func AskMnemonicPhrase() string {
	fmt.Print("Please enter your mnemonic words:")

	reader := bufio.NewReader(os.Stdin)
	mnemonic, _ := reader.ReadString('\n')

	return strings.TrimSpace(strings.ToLower(mnemonic))
}
