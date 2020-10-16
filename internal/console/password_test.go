// Copyright (c) 2020 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package console

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubPasswordReader struct {
	Passwords   []string
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
	pwdReader = &stubPasswordReader{Passwords: []string{"secret"}}
	kr = nil
	s, ok := AskPassword()
	assert.Equal(t, "secret", s)
	assert.False(t, ok)

	kr = &mockKeyring{
		p: make(map[string]string),
	}
	_ = kr.Set(service, user, "stored-in-keyring")

	pwdReader = &stubPasswordReader{Passwords: []string{"foobar"}}
	s, ok = AskPassword()
	assert.Equal(t, "stored-in-keyring", s)
	assert.True(t, ok)

	pwdReader = &stubPasswordReader{Passwords: []string{"not-stored"}}
	_ = kr.Delete(service, user)
	s, ok = AskPassword()
	assert.Equal(t, "not-stored", s)
	assert.False(t, ok)

	pwdReader = &stubPasswordReader{Passwords: []string{"not-stored"}}
	_ = StorePassword("foobar")
	s, ok = AskPassword()
	assert.Equal(t, "foobar", s)
	assert.True(t, ok)
}
