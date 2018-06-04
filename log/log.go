/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var logs *logrus.Logger

const (
	defaultLogPath       = "scan-api-log"
	defaultLogErrorLevel = "scan-api-log/error.log"
	defaultLogFile       = "all.logs"
)

//NewLogger create the log instance
func NewLogger(logLevel string, writeLog bool) *logrus.Logger {
	if logs != nil {
		return logs
	}

	logs = logrus.New()

	// get logLevel
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logs.SetLevel(level)
	}

	if writeLog {
		_ = os.Mkdir(defaultLogPath, 0777)

		path := defaultLogPath + string(os.PathSeparator) + defaultLogFile
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),       // 24 hours
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second), // 1 days
		)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}

		logs.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.InfoLevel:  writer,
				logrus.ErrorLevel: writer,
			},
			&logrus.JSONFormatter{},
		))

		pathMap := lfshook.PathMap{
			logrus.ErrorLevel: defaultLogErrorLevel,
		}
		logs.AddHook(lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{},
		))
	}

	return logs
}

// GetLogger get the default logger
func GetLogger() *logrus.Logger {
	return logs
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

//Debug output a debug message in log
func Debug(f interface{}, args ...interface{}) {
	logs.Debug(formatLog(f, args...))
}

//Info output a Info message in log
func Info(f interface{}, args ...interface{}) {
	logs.Info(formatLog(f, args...))
}

//Warn output a Warnning message in log
func Warn(f interface{}, args ...interface{}) {
	logs.Warn(formatLog(f, args...))
}

//Printf print something in console
func Printf(f interface{}, args ...interface{}) {
	logs.Print(formatLog(f, args...))
}

//Panic output a panic message in log
func Panic(f interface{}, args ...interface{}) {
	logs.Panic(formatLog(f, args...))
}

//Fatal output a fatal message in log
func Fatal(f interface{}, args ...interface{}) {
	logs.Fatal(formatLog(f, args...))
}

//Error output a error message in log
func Error(f interface{}, args ...interface{}) {
	logs.Error(formatLog(f, args...))
}
