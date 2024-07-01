package util

import (
	"encoding/json"
	"fmt"
	stdlog "log"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type MyLogrusFormatter struct {
	Fmt easy.Formatter
}

func (f *MyLogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b, e := f.Fmt.Format(entry)
	if len(entry.Data) > 0 {
		offset := 0
		strB := string(b)
		if strings.HasSuffix(strB, "\n") {
			offset = 1
		} else if strings.HasSuffix(strB, "\r\n") {
			offset = 2
		}
		jb, _ := json.Marshal(entry.Data)
		b = append(b[:len(b)-offset], []byte("\t")...)
		b = append(b, jb...)
		if offset == 1 {
			b = append(b, []byte("\n")...)
		} else if offset == 2 {
			b = append(b, []byte("\r\n")...)
		}
	}
	return b, e
}

func InitLogger(strIntLevel string) {
	logrus.SetOutput(os.Stderr)
	fmtr := MyLogrusFormatter{
		Fmt: easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n", // '%xxx%' is custom placeholder of 'logrus-easy-formatter'
		},
	}
	logrus.SetFormatter(&fmtr)
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
	if logmessage == "handle_events: error: libusb: interrupted [code -10]" {
		logrus.Debug(logmessage)
	} else {
		fmt.Printf("%s", string(data))
	}
	return len(logmessage), nil
}

func GetGoRoutineLogger(name string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"goroutine":     GoRoutineID(),
		"goroutineName": name,
	})
}
