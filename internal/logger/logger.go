package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerParams struct {
	Level   string
	LogFile string
}

var Logger *zap.SugaredLogger

func init() {
	Logger = zap.NewExample().Sugar()
}

func Init() (*zap.Logger, error) {
	var l LoggerParams
	var level zap.AtomicLevel
	switch strings.ToLower(l.Level) {
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	stdout := []string{"stdout"}
	if l.LogFile != "" {
		stdout = append(stdout, l.LogFile)
	}

	cfg := zap.Config{
		Level:       level,
		Encoding:    "console",
		OutputPaths: stdout,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}

	Logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	_ = Logger.Sync()

	return Logger, nil
}
