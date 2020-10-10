package console

import (
	"bytes"
	"errors"
	"testing"

	"github.com/chzyer/readline"
	"github.com/stretchr/testify/assert"
)

type MockBuffer struct {
	Out string
	In  []string
}

func (cb *MockBuffer) Read(p []byte) (n int, err error) {
	buf := bytes.NewBufferString(cb.In[0])
	n, err = buf.Read(p)

	if len(cb.In) > 1 {
		cb.In = cb.In[1:]
	}

	return
}

func (cb *MockBuffer) Write(p []byte) (n int, err error) {
	cb.Out += string(p)

	return len(p), nil
}

func (cb *MockBuffer) Close() error {
	return nil
}

func TestAskDoublePassword(t *testing.T) {
	mb := &MockBuffer{In: []string{"secret\n", "secret\n"}}
	readliner, _ = readline.NewEx(&readline.Config{
		Stdin:  mb,
		Stdout: mb,
	})

	b, err := AskDoublePassword()
	assert.NoError(t, err)
	assert.Equal(t, "secret", string(b))
	assert.Contains(t, mb.Out, "Please enter your vault password")
	assert.Contains(t, mb.Out, "Please retype your vault password")


	mb = &MockBuffer{In: []string{"secret\n", "false\n", "secret\n", "secret\n"}}
	readliner, _ = readline.NewEx(&readline.Config{
		Stdin:  mb,
		Stdout: mb,
	})

	b, err = AskDoublePassword()
	assert.NoError(t, err)
	assert.Equal(t, "secret", string(b))
	assert.Contains(t, mb.Out, "Please enter your vault password")
	assert.Contains(t, mb.Out, "Please retype your vault password")
	assert.Contains(t, mb.Out, "Passwords do not match. Please type again.")
}

type mockKeyring struct {
	p map[string]string
}

func (m *mockKeyring) Set(service, user, password string) error {
	m.p[service+":"+user] = password
	return nil
}

func (m *mockKeyring) Get(service, user string) (string, error) {
	v, ok := m.p[service+":"+user]
	if !ok {
		return "", errors.New("not found")
	}

	return v, nil
}

func (m *mockKeyring) Delete(service, user string) error {
	delete(m.p, service+":"+user)
	return nil
}

func TestAskPassword(t *testing.T) {
	mb := &MockBuffer{In: []string{"secret\n"}}
	readliner, _ = readline.NewEx(&readline.Config{
		Stdin:  mb,
		Stdout: mb,
	})

	kr = nil
	s, ok := AskPassword()
	assert.Equal(t, "secret", s)
	assert.False(t, ok)

	kr = &mockKeyring{
		p: make(map[string]string, 0),
	}
	_ = kr.Set(service, user, "stored-in-keyring")


	mb = &MockBuffer{In: []string{"not-stored\n"}}
	readliner, _ = readline.NewEx(&readline.Config{
		Stdin:  mb,
		Stdout: mb,
	})

	s, ok = AskPassword()
	assert.Equal(t, "stored-in-keyring", s)
	assert.True(t, ok)

	_ = kr.Delete(service, user)
	s, ok = AskPassword()
	assert.Equal(t, "not-stored", s)
	assert.False(t, ok)

	_ = StorePassword("foobar")
	s, ok = AskPassword()
	assert.Equal(t, "foobar", s)
	assert.True(t, ok)
}
