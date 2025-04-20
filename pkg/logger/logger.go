package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Setup инициализирует логгер
func Setup() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

// GetLogger возвращает экземпляр логгера
func GetLogger() *zap.Logger {
	return log
}

// Info логирует информационное сообщение
func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

// Error логирует сообщение об ошибке
func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}

// Debug логирует отладочное сообщение
func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

// Warn логирует предупреждение
func Warn(message string, fields ...zap.Field) {
	log.Warn(message, fields...)
}

// Fatal логирует фатальную ошибку и завершает программу
func Fatal(message string, fields ...zap.Field) {
	log.Fatal(message, fields...)
}
