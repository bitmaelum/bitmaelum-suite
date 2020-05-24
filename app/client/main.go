package main

import (
	"github.com/jaytaph/mailv2/cmd"
	"github.com/sirupsen/logrus"
	logger "github.com/sirupsen/logrus"
	"os"
)

func main() {
	// This should probably be in the root of command so we can deal with flags like -vvv etc
	logger.SetFormatter(new(logrus.JSONFormatter))
	logger.SetFormatter(new(logrus.TextFormatter))
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetOutput(os.Stdout)

	cmd.Execute()
}
