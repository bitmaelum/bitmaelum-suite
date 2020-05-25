package client

import (
	"github.com/jaytaph/mailv2/client/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	// This should probably be in the root of command so we can deal with flags like -vvv etc
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetOutput(os.Stdout)

	cmd.Execute()
}
