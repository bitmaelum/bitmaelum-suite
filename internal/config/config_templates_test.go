// Copyright (c) 2021 BitMaelum Authors
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

package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplates(t *testing.T) {
	var buf = bytes.Buffer{}
	err := GenerateServerConfig(&buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.String())

	assert.Empty(t, Server.Logging.Level)
	err = Server.LoadConfig(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "trace", Server.Logging.Level)

	err = GenerateClientConfig(&buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.String())

	assert.Empty(t, Client.Resolvers.Remote.URL)
	err = Client.LoadConfig(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "https://resolver.bitmaelum.com", Client.Resolvers.Remote.URL)

	err = GenerateBridgeConfig(&buf)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.String())

	assert.Empty(t, Bridge.Server.SMTP.Host)
	err = Bridge.LoadConfig(&buf)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", Bridge.Server.SMTP.Host)

}
