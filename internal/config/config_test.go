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
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	hook *test.Hook
)

func TestClientConfig(t *testing.T) {
	testingClientConfigPath := "/etc/bitmaelum/bitmaelum-client-config.yml"
	testingNotExistsFilePath := "/etc/bitmaelum/not-exist.yml"
	testingResolverRemoteURL := "https://resolver.bitmaelum.com"

	_ = os.Setenv("BITMAELUM_CLIENT_CONFIG", "")

	fs = afero.NewMemMapFs()

	err := LoadClientConfigOrPass("")
	assert.Error(t, err)

	f, err := fs.Create(testingClientConfigPath)
	assert.NoError(t, err)
	err = GenerateClientConfig(f)
	assert.NoError(t, err)
	_ = f.Close()

	Client.Resolver.Remote.URL = ""
	err = LoadClientConfigOrPass(testingClientConfigPath)
	assert.NoError(t, err)
	assert.Equal(t, testingResolverRemoteURL, Client.Resolver.Remote.URL)

	Client.Resolver.Remote.URL = ""
	err = LoadClientConfigOrPass(testingNotExistsFilePath)
	assert.Error(t, err)
	assert.Equal(t, "", Client.Resolver.Remote.URL)

	Client.Resolver.Remote.URL = ""
	err = LoadClientConfigOrPass(testingClientConfigPath)
	assert.NoError(t, err)
	assert.Equal(t, testingResolverRemoteURL, Client.Resolver.Remote.URL)

	// Read from searchpath
	if runtime.GOOS != "windows" {
		// This test fails on windows. It expects the file to be on the searchpath, but it isn't because
		// the searchpath for windows is different. However, we expect a regular path like /etc/bitmaelum/*.yml
		// for this test to succeed.
		Client.Resolver.Remote.URL = ""
		err = LoadClientConfigOrPass("")
		assert.NoError(t, err)
		assert.Equal(t, testingResolverRemoteURL, Client.Resolver.Remote.URL)
	}

	// Read from non-existing env
	Client.Resolver.Remote.URL = ""
	_ = os.Setenv("BITMAELUM_CLIENT_CONFIG", "/etc/does/not/exist.yml")
	err = LoadClientConfigOrPass("")
	assert.Error(t, err)

	// Read from existing env
	Client.Resolver.Remote.URL = ""
	_ = os.Setenv("BITMAELUM_CLIENT_CONFIG", testingClientConfigPath)
	err = LoadClientConfigOrPass("")
	assert.NoError(t, err)
	assert.Equal(t, testingResolverRemoteURL, Client.Resolver.Remote.URL)
}

func TestServerConfig(t *testing.T) {
	testingNotExistsFilePath := "/etc/bitmaelum/not-exist.yml"
	testingServerConfigPath := "/etc/bitmaelum/bitmaelum-server-config.yml"

	_ = os.Setenv("BITMAELUM_SERVER_CONFIG", "")

	err := LoadServerConfigOrPass("")
	assert.Error(t, err)

	fs = afero.NewMemMapFs()
	f, err := fs.Create(testingServerConfigPath)
	assert.NoError(t, err)
	err = GenerateServerConfig(f)
	assert.NoError(t, err)
	_ = f.Close()

	// Load direct
	Server.Work.Pow.Bits = 0
	err = LoadServerConfigOrPass(testingServerConfigPath)
	assert.NoError(t, err)
	assert.Equal(t, 25, Server.Work.Pow.Bits)

	// Unknown file
	Server.Work.Pow.Bits = 0
	err = LoadServerConfigOrPass(testingNotExistsFilePath)
	assert.Error(t, err)
	assert.Equal(t, 0, Server.Work.Pow.Bits)

	// Load direct
	Server.Work.Pow.Bits = 0
	err = LoadServerConfigOrPass(testingServerConfigPath)
	assert.NoError(t, err)
	assert.Equal(t, 25, Server.Work.Pow.Bits)

	// Read from predetermined paths
	if runtime.GOOS != "windows" {
		// This test fails on windows. It expects the file to be on the searchpath, but it isn't because
		// the searchpath for windows is different. However, we expect a regular path like /etc/bitmaelum/*.yml
		// for this test to succeed.

		Server.Work.Pow.Bits = 0
		err = LoadServerConfigOrPass("")
		assert.NoError(t, err)
		assert.Equal(t, 25, Server.Work.Pow.Bits)
	}

	// Read from non-existing env
	Server.Work.Pow.Bits = 0
	_ = os.Setenv("BITMAELUM_SERVER_CONFIG", "/etc/does/not/exist.yml")
	err = LoadServerConfigOrPass("")
	assert.Error(t, err)

	// Read from existing env
	Server.Work.Pow.Bits = 0
	_ = os.Setenv("BITMAELUM_SERVER_CONFIG", testingServerConfigPath)
	err = LoadServerConfigOrPass("")
	assert.NoError(t, err)
	assert.Equal(t, 25, Server.Work.Pow.Bits)
}

func TestLoadClientConfig(t *testing.T) {
	// Failed loading
	err := readConfigPath("/foo/bar", "", Client.LoadConfig, &LoadedClientConfigPath)
	assert.Error(t, err)
}

func TestGenerateRoutingFromMnemonic(t *testing.T) {
	r, err := GenerateRoutingFromMnemonic("ed25519 cluster puppy wash ceiling skate search great angry drift rose undo fragile boring fence stumble shuffle cable praise")
	assert.NoError(t, err)

	assert.Equal(t, "f5f1dc4eff7237ac0e061a9e8982b7b913fc479138189cc8d6ba5131dee1bde9", r.RoutingID)
	assert.Equal(t, "ed25519 MC4CAQAwBQYDK2VwBCIEIDLOvf5iUAPWeNIYlbyDffgv+VA2xnS1s1mUYIOmW8XK", r.KeyPair.PrivKey.String())
	assert.Equal(t, "ed25519 MCowBQYDK2VwAyEAndS2/G3uasbaYO0+89rNzvNJ3gfOi/An1t5xvETeNoc=", r.KeyPair.PubKey.String())
}

func init() {
	// Setup mock
	_, hook = test.NewNullLogger()
	logrus.AddHook(hook)
	logrus.SetOutput(ioutil.Discard)
	logrus.StandardLogger().ExitFunc = func(int) {
		// dummy function to prevent os.Exit being called by logrus
	}
}
