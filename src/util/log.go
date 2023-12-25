package util

import (
	stdlog "log"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func InitLogger(strIntLevel string) {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})
	SetLogLevel(strIntLevel)
	stdlog.SetOutput(new(LogrusWriter))
}

func SetLogLevel(strIntLevel string) {
	loglevel := logrus.InfoLevel
	switch strings.ToLower(strIntLevel) {
	case "panic":
		loglevel = logrus.PanicLevel
	case "fatal":
		loglevel = logrus.FatalLevel
	case "error":
		loglevel = logrus.ErrorLevel
	case "warn":
		loglevel = logrus.WarnLevel
	case "info":
		loglevel = logrus.InfoLevel
	case "debug":
		loglevel = logrus.DebugLevel
	case "trace":
		loglevel = logrus.TraceLevel
	}
	logrus.SetLevel(loglevel)
}

/********** work around to this problem **********/
// https://github.com/google/gousb/issues/87#issuecomment-1100956460
type LogrusWriter int

var reStdGoLogFormat = regexp.MustCompile(`(?s)[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} (?P<msg>.+)`)

func (LogrusWriter) Write(data []byte) (int, error) {
	logmessage := string(data)
	if reStdGoLogFormat.MatchString(logmessage) {
		logmessage = logmessage[20:]
	}
	if strings.HasSuffix(logmessage, "\n") {
		logmessage = logmessage[:len(logmessage)-1]
	}
	logrus.Infof("[gousb] %s", logmessage)
	return len(logmessage), nil
}
