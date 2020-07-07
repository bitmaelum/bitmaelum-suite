package config

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var (
	fatal bool = false
	hook  *test.Hook
)

func TestLoadClientConfig(t *testing.T) {
	// Failed loading
	LoadClientConfig("/foo/bar")
	assert.True(t, fatal)
	assert.Len(t, hook.Entries, 1)
	assert.Equal(t, "cannot load client configuration", hook.Entries[0].Message)
}

func init() {
	// Setup mock
	_, hook = test.NewNullLogger()
	logrus.AddHook(hook)
	logrus.SetOutput(ioutil.Discard)
	logrus.StandardLogger().ExitFunc = func(int) { fatal = true }
}
