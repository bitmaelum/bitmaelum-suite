package internal

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var fs = afero.NewOsFs()

// SetLogging will set the correct level and log path
func SetLogging(level, path string) {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = time.Stamp
	logrus.SetFormatter(formatter)

	// Default to stderr
	logrus.SetOutput(os.Stderr)

	if path == "stdout" {
		logrus.SetOutput(os.Stdout)

	} else if path == "stderr" {
		logrus.SetOutput(os.Stderr)

	} else if strings.HasPrefix(path, "syslog") {
		// Default to localhost syslog daemon
		syslogHost := "localhost:514"

		splits := strings.SplitN(path, ":", 2)
		if len(splits) == 2 {
			syslogHost = splits[1]
		}

		hook, err := setupSyslogHook("udp", syslogHost)
		if err != nil {
			logrus.Error("Unable to connect to syslog daemon. Falling back to stderr")
			logrus.SetOutput(os.Stderr)
		} else {
			logrus.AddHook(hook)
		}
	} else {
		// Default to a path
		w, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			w = os.Stderr
		}

		logrus.SetOutput(w)
	}

	switch level {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}

	logrus.Tracef("setting loglevel to '%s'", level)
}
