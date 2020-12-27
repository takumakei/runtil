package app

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZap(level zapcore.Level) error {
	logger, err := initZap(level)
	if err == nil {
		zap.ReplaceGlobals(logger)
	}
	return err
}

func initZap(level zapcore.Level) (*zap.Logger, error) {
	if FlagDebug() {
		config := zap.NewDevelopmentConfig()
		config.Level.SetLevel(level)
		return config.Build()
	}
	config := zap.NewProductionConfig()
	config.Level.SetLevel(level)
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	return config.Build()
}
