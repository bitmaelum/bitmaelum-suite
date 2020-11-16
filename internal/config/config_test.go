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

package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	fatal = false
	hook  *test.Hook
)

func TestClientConfig(t *testing.T) {
	fs = afero.NewMemMapFs()

	err := LoadClientConfigOrPass("")
	assert.Error(t, err)

	f, err := fs.Create("/etc/bitmaelum/client-config.yml")
	assert.NoError(t, err)
	err = GenerateClientConfig(f)
	assert.NoError(t, err)
	_ = f.Close()

	Client.Accounts.ProofOfWork = 0
	err = LoadClientConfigOrPass("/etc/bitmaelum/client-config.yml")
	assert.NoError(t, err)
	assert.Equal(t, 22, Client.Accounts.ProofOfWork)

	Client.Accounts.ProofOfWork = 0
	err = LoadClientConfigOrPass("/etc/bitmaelum/not-exist.yml")
	assert.Error(t, err)
	assert.Equal(t, 0, Client.Accounts.ProofOfWork)

	Client.Accounts.ProofOfWork = 0
	err = LoadClientConfigOrPass("/etc/bitmaelum/client-config.yml")
	assert.NoError(t, err)
	assert.Equal(t, 22, Client.Accounts.ProofOfWork)

	Client.Accounts.ProofOfWork = 0
	err = LoadClientConfigOrPass("")
	assert.Error(t, err)
	assert.Equal(t, 0, Client.Accounts.ProofOfWork)

	// Read from non-existing env
	Client.Accounts.ProofOfWork = 0
	os.Setenv("BITMAELUM_CLIENT_CONFIG", "/etc/does/not/exist.yml")
	err = LoadClientConfigOrPass("")
	assert.Error(t, err)

	// Read from existing env
	Client.Accounts.ProofOfWork = 0
	os.Setenv("BITMAELUM_CLIENT_CONFIG", "/etc/bitmaelum/client-config.yml")
	err = LoadClientConfigOrPass("")
	assert.NoError(t, err)
	assert.Equal(t, 22, Client.Accounts.ProofOfWork)
}

func TestServerConfig(t *testing.T) {
	err := LoadServerConfigOrPass("")
	assert.Error(t, err)

	fs = afero.NewMemMapFs()
	f, err := fs.Create("/etc/bitmaelum/server-config.yml")
	assert.NoError(t, err)
	err = GenerateClientConfig(f)
	assert.NoError(t, err)
	_ = f.Close()

	// Load direct
	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("/etc/bitmaelum/server-config.yml")
	assert.NoError(t, err)
	assert.Equal(t, 22, Server.Accounts.ProofOfWork)

	// Unknown file
	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("/etc/bitmaelum/not-exist.yml")
	assert.Error(t, err)
	assert.Equal(t, 0, Server.Accounts.ProofOfWork)

	// Load direct
	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("/etc/bitmaelum/server-config.yml")
	assert.NoError(t, err)
	assert.Equal(t, 22, Server.Accounts.ProofOfWork)

	// Read from predetermined paths
	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("")
	assert.Error(t, err)
	assert.Equal(t, 0, Server.Accounts.ProofOfWork)

	// Read from non-existing env
	Server.Accounts.ProofOfWork = 0
	os.Setenv("BITMAELUM_SERVER_CONFIG", "/etc/does/not/exist.yml")
	err = LoadServerConfigOrPass("")
	assert.Error(t, err)

	// Read from existing env
	Server.Accounts.ProofOfWork = 0
	os.Setenv("BITMAELUM_SERVER_CONFIG", "/etc/bitmaelum/server-config.yml")
	err = LoadServerConfigOrPass("")
	assert.NoError(t, err)
	assert.Equal(t, 22, Server.Accounts.ProofOfWork)
}

func TestLoadClientConfig(t *testing.T) {
	// Failed loading
	err := readConfigPath("/foo/bar", "", Client.LoadConfig)
	assert.Error(t, err)
}

func TestGenerateRoutingFromMnemonic(t *testing.T) {
	r, err := GenerateRoutingFromMnemonic("cluster puppy wash ceiling skate search great angry drift rose undo fragile boring fence stumble shuffle cable praise")
	assert.NoError(t, err)

	assert.Equal(t, "f5f1dc4eff7237ac0e061a9e8982b7b913fc479138189cc8d6ba5131dee1bde9", r.RoutingID)
	assert.Equal(t, "ed25519 MC4CAQAwBQYDK2VwBCIEIDLOvf5iUAPWeNIYlbyDffgv+VA2xnS1s1mUYIOmW8XK", r.PrivateKey.String())
	assert.Equal(t, "ed25519 MCowBQYDK2VwAyEAndS2/G3uasbaYO0+89rNzvNJ3gfOi/An1t5xvETeNoc=", r.PublicKey.String())
}

func init() {
	// Setup mock
	_, hook = test.NewNullLogger()
	logrus.AddHook(hook)
	logrus.SetOutput(ioutil.Discard)
	logrus.StandardLogger().ExitFunc = func(int) { fatal = true }
}
