package console

import (
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

// AutoComplete Listener
type acl struct{}

func (l *acl) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// Empty string will return empty
	s := string(line)[:pos]
	if s == "" {
		return
	}

	// Find word with the current prefix
	var f *string
	for _, w := range bip39.GetWordList() {
		if strings.HasPrefix(w, s) {
			f = &w
			break
		}
	}

	// Can't find the word, redo the word
	if f == nil {
		return line[:pos-1], pos - 1, true
	}

	return []rune(*f), pos, true
}

// AskMnemonicPhrase asks for the mnemonic phrase, which is autocompleted.
func AskMnemonicPhrase() string {
	fmt.Print("Please enter your mnemonic words and press enter after each word. End with an empty line.")

	var s = ""
	for {
		line, err := readliner.Readline()
		if err != nil {
			return ""
		}

		if line == "" {
			break
		}

		s += line + " "
	}

	return strings.TrimSpace(s)
}
