package config

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
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

	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("/etc/bitmaelum/server-config.yml")
	assert.NoError(t, err)
	assert.Equal(t, 22, Server.Accounts.ProofOfWork)

	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("/etc/bitmaelum/not-exist.yml")
	assert.Error(t, err)
	assert.Equal(t, 0, Server.Accounts.ProofOfWork)

	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("/etc/bitmaelum/server-config.yml")
	assert.NoError(t, err)
	assert.Equal(t, 22, Server.Accounts.ProofOfWork)

	Server.Accounts.ProofOfWork = 0
	err = LoadServerConfigOrPass("")
	assert.NoError(t, err)
	assert.Equal(t, 22, Server.Accounts.ProofOfWork)
}

func TestLoadClientConfig(t *testing.T) {
	// Failed loading
	err := readConfigPath("/foo/bar", Client.LoadConfig)
	assert.Error(t, err)
}

func init() {
	// Setup mock
	_, hook = test.NewNullLogger()
	logrus.AddHook(hook)
	logrus.SetOutput(ioutil.Discard)
	logrus.StandardLogger().ExitFunc = func(int) { fatal = true }
}
