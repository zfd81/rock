package otto

import (
	log "github.com/sirupsen/logrus"
)

func LogInfo(args ...interface{}) {
	log.Info(args...)
}

func LogError(args ...interface{}) {
	log.Error(args...)
}
