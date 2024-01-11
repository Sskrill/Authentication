package transport

import "github.com/sirupsen/logrus"

func loggerFields(handler string) logrus.Fields {
	return logrus.Fields{"handler": handler}
}
func loggerError(handler string, err error) {
	logrus.WithFields(loggerFields(handler)).Error(err)
}
