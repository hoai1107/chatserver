package logwrapper

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var once sync.Once

func InitLogger() {
	once.Do(func() {
		log = &logrus.Logger{
			Out:       os.Stdout,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
		}

		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	})
}

func Info(args ...interface{}) {
	log.Infoln(args...)
}

// FIXME: Circular import
// func InfoMessage(msg chat.Message) {
// 	log.WithFields(logrus.Fields{
// 		"type":     msg.Type,
// 		"username": msg.Username,
// 		"content":  msg.Content,
// 	}).Info("Get message.")
// }

func Error(args ...interface{}) {
	log.Errorln(args...)
}
