package logger

import (
	"go-rest-api/config"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init(cfg config.Configuration) {
	Logger = newLogger(cfg)
}

func newLogger(cfg config.Configuration) *zap.SugaredLogger {
	var logger *zap.Logger
	if cfg.LoggerLevel == "dev" {
		logger, _ = zap.NewDevelopment()
	} else if cfg.LoggerLevel == "prodaction" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	sugaredLogger := logger.Sugar()
	return sugaredLogger
}
