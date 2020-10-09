package console

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/tyler-smith/go-bip39"
)

// AutoComplete Listener
type acl struct{}

func (l *acl) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	s := string(line)[:pos]
	if s == "" {
		return
	}

	var f *string
	for _, w := range bip39.GetWordList() {
		if strings.HasPrefix(w, s) {
			f = &w
			break
		}
	}

	if f == nil {
		return line[:pos-1], pos - 1, true
	}

	return []rune(*f), pos, true
}

// AskSeedPhrase asks for the seed phrase, which is autocompleted.
func AskSeedPhrase() string {
	fmt.Print("Please enter your seed words: ")

	rl, _ := readline.NewEx(&readline.Config{
		Prompt:   ">>> ",
		Listener: &acl{},
	})
	defer func() {
		_ = rl.Close()
	}()

	var s = ""
	for {
		line, err := rl.Readline()
		if err != nil {
			return ""
		}

		if line == "" {
			break
		}

		s += line + " "
	}

	return s
}
