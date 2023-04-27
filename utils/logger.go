package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/sirupsen/logrus"
)

type Logrus struct {
	TimestampFormat string
	LogFormat       string
}

func (f *Logrus) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = constants.DefaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = constants.DefaultTimestampFormat
	}

	var levelColor int
	switch entry.Level {
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31
	case logrus.WarnLevel:
		levelColor = 33
	case logrus.DebugLevel:
		levelColor = 34
	case logrus.TraceLevel:
		levelColor = 35
	default:
		levelColor = 36
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	level := strings.ToUpper(entry.Level.String())
	colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, level)
	output = strings.Replace(output, "%lvl%", colored, 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}

	return []byte(output), nil
}
