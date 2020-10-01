// +build !windows

package internal

import (
	"log/syslog"

	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

func setupSyslogHook(proto, host string) (logrus.Hook, error) {
	return logrus_syslog.NewSyslogHook(proto, host, syslog.LOG_DAEMON, "BitMaelum")
}
