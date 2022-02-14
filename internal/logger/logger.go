package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	lgr *zap.Logger
}

func getLogLevel(level string) zapcore.Level {
	var logLevel zapcore.Level
	switch level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	return logLevel
}

func buildConfig(level string, outputs []string) (*zap.Logger, error) {
	config := &zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(getLogLevel(level)),
		OutputPaths:      outputs,
		ErrorOutputPaths: outputs,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "msg",
			LevelKey:      "lvl",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "time",
			EncodeTime:    zapcore.RFC3339TimeEncoder,
			CallerKey:     "",
			EncodeCaller:  zapcore.ShortCallerEncoder,
			StacktraceKey: "",
		},
	}

	return config.Build()
}

func NewLogger(level string, outputs []string) (*Logger, error) {
	lgr, err := buildConfig(level, outputs)

	return &Logger{lgr: lgr}, err
}

func (l *Logger) Debug(msg string, values ...interface{}) {
	l.lgr.Sugar().Debugw(msg, values...)
}

func (l *Logger) Info(msg string, values ...interface{}) {
	l.lgr.Sugar().Infow(msg, values...)
}

func (l *Logger) Warn(msg string, values ...interface{}) {
	l.lgr.Sugar().Warnw(msg, values...)
}

func (l *Logger) Error(msg string, values ...interface{}) {
	l.lgr.Sugar().Errorw(msg, values...)
}

func (l *Logger) Lgr() *zap.Logger {
	return l.lgr
}
