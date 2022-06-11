package logwrapper

import (
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var once sync.Once

func InitLogger() {
	once.Do(func() {
		log = logrus.New()

		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	})
}

func Info(args ...interface{}) {
	log.Infoln(args...)
}

func InfoMessage(msg interface{}) {
	val := reflect.ValueOf(msg).Elem()
	msgType := val.FieldByName("Type").Interface().(string)
	msgUsername := val.FieldByName("Username").Interface().(string)
	msgContent := val.FieldByName("Content").Interface().(string)

	log.WithFields(logrus.Fields{
		"type":     msgType,
		"username": msgUsername,
		"content":  msgContent,
	}).Info("Get message.")
}

func Error(args ...interface{}) {
	log.Errorln(args...)
}
