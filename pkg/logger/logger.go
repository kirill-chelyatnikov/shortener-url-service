package logger

import (
	"go.uber.org/zap"
)

// InitLogger - инициализация логгера
func InitLogger() *zap.SugaredLogger {
	log := zap.Must(zap.NewProduction())
	defer log.Sync()

	return log.Sugar()
}
