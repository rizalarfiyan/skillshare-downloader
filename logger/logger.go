package logger

import (
	"os"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/utils"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func Init() {
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.InfoLevel,
		Formatter: &utils.Logrus{
			TimestampFormat: constants.DefaultTimestampFormat,
			LogFormat:       constants.DefaultLogFormat,
		},
	}
}

func Get() *logrus.Logger {
	return logger
}

func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}
