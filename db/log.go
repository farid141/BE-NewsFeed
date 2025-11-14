package db

import (
	"os"

	"github.com/farid141/go-rest-api/config"
	"github.com/sirupsen/logrus"
)

func NewLog(cfg config.Config) *logrus.Logger {
	logger := logrus.New()

	file, _ := os.OpenFile(cfg.LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logger.SetOutput(file)

	return logger
}
