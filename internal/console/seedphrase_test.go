package console

import (
	"testing"

	"github.com/chzyer/readline"
	"github.com/stretchr/testify/assert"
)

func TestAcl_OnChange(t *testing.T) {

	l := acl{}

	var (
		nl []rune
		np int
		ok bool
	)

	nl, np, ok = l.OnChange([]rune(""), 0, 'a')
	assert.Equal(t, "", string(nl))
	assert.Equal(t, 0, np)
	assert.False(t, ok)

	nl, np, ok = l.OnChange([]rune("s"), 1, 'a')
	assert.Equal(t, "sad", string(nl))
	assert.Equal(t, 1, np)
	assert.True(t, ok)

	nl, np, ok = l.OnChange([]rune("ex"), 2, 'x')
	assert.Equal(t, "exact", string(nl))
	assert.Equal(t, 2, np)
	assert.True(t, ok)

	nl, np, ok = l.OnChange([]rune("ep"), 2, ' ')
	assert.Equal(t, "episode", string(nl))
	assert.Equal(t, 2, np)
	assert.True(t, ok)

	nl, np, ok = l.OnChange([]rune("asf"), 3, ' ')
	assert.Equal(t, "as", string(nl))
	assert.Equal(t, 2, np)
	assert.True(t, ok)
}

func TestAskSeedPhrase(t *testing.T) {
	mb := &MockBuffer{In: []string{"foo\n", "bar\n", "baz\n", "\n"}}
	readliner, _ = readline.NewEx(&readline.Config{
		Stdin:  mb,
		Stdout: mb,
		Stderr: mb,
	})

	s := AskSeedPhrase()
	assert.Equal(t, "foo bar baz", s)
}
