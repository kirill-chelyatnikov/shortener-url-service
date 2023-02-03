package logger

import "github.com/sirupsen/logrus"

// InitLogger - инициализация логгера
func InitLogger() *logrus.Logger {
	log := logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	log.SetFormatter(formatter)

	return log
}
