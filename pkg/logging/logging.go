package logging

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

// Init initializes the logger
func Init() (err error) {
	conf := zap.NewDevelopmentConfig()
	//conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err = conf.Build()
	return
}

// Logger returns our global logger
func Logger() *zap.Logger {
	if logger == nil {

	}
	return logger
}
