// +build !windows

package core

import (
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

func setupSyslogHook(proto, host string) (logrus.Hook, error) {
	return logrus_syslog.NewSyslogHook(proto, host, syslog.LOG_DAEMON, "BitMaelum")
}
