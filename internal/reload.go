package internal

import "github.com/sirupsen/logrus"

// Reload will reload any configuration changes without restarting the server
func Reload() {
	logrus.Info("Reloading configurations (not implemented)")
}
