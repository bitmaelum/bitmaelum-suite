package config

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var (
	fatal = false
	hook  *test.Hook
)

func TestLoadClientConfig(t *testing.T) {
	// Failed loading
	err := readConfigPath("/foo/bar", Client.LoadConfig)
	assert.EqualError(t, err, "open /foo/bar: no such file or directory")
}

func init() {
	// Setup mock
	_, hook = test.NewNullLogger()
	logrus.AddHook(hook)
	logrus.SetOutput(ioutil.Discard)
	logrus.StandardLogger().ExitFunc = func(int) { fatal = true }
}
