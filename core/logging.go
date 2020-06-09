package core

import (
    "github.com/jaytaph/mailv2/core/config"
    "github.com/sirupsen/logrus"
    "os"
)

func SetLogging(level string) {
    // How do we know how we need to log?
    logrus.SetFormatter(new(logrus.JSONFormatter))
    logrus.SetFormatter(new(logrus.TextFormatter))

    // We probably want to set this as well through params
    logrus.SetOutput(os.Stdout)


    switch (level) {
    case "trace":
        logrus.SetLevel(logrus.TraceLevel)
        break;
    case "debug":
        logrus.SetLevel(logrus.DebugLevel)
        break;
    case "info":
        logrus.SetLevel(logrus.InfoLevel)
        break;
    case "warning":
        logrus.SetLevel(logrus.WarnLevel)
        break;
    case "error":
    default:
        logrus.SetLevel(logrus.ErrorLevel)
        break;
    }

    logrus.Tracef("setting loglevel to '%s'", config.Server.Logging.Level)
}
