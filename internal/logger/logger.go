package logger

import (
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Params struct {
	Level   string `json:"level" mapstructure:"level"`
	LogFile string `json:"file" mapstructure:"file"`
}

var Logger *zap.SugaredLogger

func init() {
	Logger = zap.NewExample().Sugar()
}

func Init(l Params) (*zap.SugaredLogger, error) {
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

	outputPath := []string{"stdout"}
	errorOutputPath := []string{"stderr"}

	if l.LogFile != "" {
		_, err := os.Create(l.LogFile)
		if err != nil {
			log.Printf("Can't create config file %s %v", l.LogFile, err)
		} else {
			outputPath = append(outputPath, l.LogFile)
			errorOutputPath = append(errorOutputPath, l.LogFile)
		}
	}

	cfg := zap.Config{
		Level:            level,
		Encoding:         "json",
		OutputPaths:      outputPath,
		ErrorOutputPaths: errorOutputPath,
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

	return Logger.Sugar(), nil
}
