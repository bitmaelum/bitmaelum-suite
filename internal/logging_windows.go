// +build windows

package internal

import (
	"errors"
	"github.com/sirupsen/logrus"
)

func setupSyslogHook(proto, host string) (logrus.Hook, error) {
	return nil, errors.New("syslog not implemented on windows")
}
