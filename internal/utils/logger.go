package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Logger *zap.Logger
	LogsLevelController *zap.AtomicLevel
}

func NewLogger(logsLevel string) *Logger {
	var level zapcore.Level

	switch logsLevel {
	case "debug":
		level = zap.DebugLevel
	case "error":
		level = zap.ErrorLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	default:
		level = zap.InfoLevel
	}

	levelController := zap.NewAtomicLevelAt(level)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			levelController,
		),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return &Logger{
		Logger: logger,
		LogsLevelController: &levelController,
	}
}

func (l *Logger) SetLogsLevel(level string) {
	var logsLevel zapcore.Level

	switch level {
	case "debug":
		logsLevel = zap.DebugLevel
	case "error":
		logsLevel = zap.ErrorLevel
	case "info":
		logsLevel = zap.InfoLevel
	case "warn":
		logsLevel = zap.WarnLevel
	default:
		logsLevel = zap.InfoLevel
	}

	l.LogsLevelController.SetLevel(logsLevel)
}