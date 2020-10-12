package console

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubPasswordReader struct {
    Passwords  []string
    ReturnError bool
}

func (pr *stubPasswordReader) ReadPassword() ([]byte, error) {
    if pr.ReturnError {
        return nil, errors.New("stubbed error")
    }

    s := pr.Passwords[0]
	pr.Passwords = pr.Passwords[1:]

    return []byte(s), nil
}


func TestAskDoublePassword(t *testing.T) {
	pwdReader = &stubPasswordReader{ReturnError: true}
	_, err := AskDoublePassword()
	assert.Error(t, err)

	pwdReader = &stubPasswordReader{Passwords: []string{"secret", "secret"}}
	b, err := AskDoublePassword()
	assert.NoError(t, err)
	assert.Equal(t, "secret", string(b))

	pwdReader = &stubPasswordReader{Passwords: []string{"secret", "foo", "bar", "bar"}}
	b, err = AskDoublePassword()
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(b))
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
	kr = nil
	s, ok := AskPassword()
	assert.Equal(t, "secret", s)
	assert.False(t, ok)

	kr = &mockKeyring{
		p: make(map[string]string),
	}
	_ = kr.Set(service, user, "stored-in-keyring")

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
