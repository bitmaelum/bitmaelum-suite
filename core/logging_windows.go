// +build windows

package core

import "github.com/sirupsen/logrus"

func setupSyslogHook(proto, host string) (logrus.Hook, err) {
	return nil, errors.New("syslog not implemented on windows")
}
