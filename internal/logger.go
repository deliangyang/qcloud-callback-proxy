package internal

import (
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	al := &logrus.Logger{}

	log = al
}

// Logger logger entity
func Logger() *logrus.Logger {
	return log
}
