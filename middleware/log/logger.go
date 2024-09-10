package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func initialLogger() *zap.Logger {

	encoderCfg := zap.NewProductionEncoderConfig()
	leves := zap.NewAtomicLevelAt(zap.InfoLevel)
	if os.Getenv("APP_ENV") == "development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
		leves = zap.NewAtomicLevelAt(zap.DebugLevel)

	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             leves,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		EncoderConfig:     encoderCfg,
		Encoding:          "json",
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

func LoadLogger() {
	Logger = initialLogger()
}
