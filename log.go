package asyncq

import (
	"github.com/sirupsen/logrus"
)

var log logger = logrus.New()

func SetLogger(l logger) {
	log = l
}

type logger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}
