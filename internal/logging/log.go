package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"runtime"
	"strings"
)

var (
	w *logWrapper
)

type logWrapper struct {
	reportCaller bool
	logger       *log.Logger
}

func init() {
	w = &logWrapper{
		reportCaller: true,
		logger:       log.New(),
	}
}

func Setup(level log.Level, output io.Writer, formatter log.Formatter, reportCaller bool) {
	w.logger.SetFormatter(formatter)
	w.logger.SetOutput(output)
	w.logger.SetLevel(level)
	w.reportCaller = reportCaller
}

// Panic logs a message at PanicLevel.
func Panic(args ...interface{}) {
	logWithFields(log.PanicLevel, log.Fields{}, args...)
}

// Fatal logs a message at FatalLevel.
func Fatal(args ...interface{}) {
	logWithFields(log.FatalLevel, log.Fields{}, args...)
}

// Error logs a message at ErrorLevel.
func Error(args ...interface{}) {
	logWithFields(log.ErrorLevel, log.Fields{}, args...)
}

// Warn logs a message at WarnLevel.
func Warning(args ...interface{}) {
	logWithFields(log.WarnLevel, log.Fields{}, args...)
}

// Info logs a message at InfoLevel.
func Info(args ...interface{}) {
	logWithFields(log.InfoLevel, log.Fields{}, args...)
}

// Debug logs a message at DebugLevel.
func Debug(args ...interface{}) {
	logWithFields(log.DebugLevel, log.Fields{}, args...)
}

// Trace logs a message at TraceLevel.
func Trace(args ...interface{}) {
	logWithFields(log.TraceLevel, log.Fields{}, args...)
}

// PanicWithFields logs a message at PanicLevel. Including any fields passed
func PanicWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(log.PanicLevel, fields, args...)
}

// FatalWithFields logs a message at FatalLevel. Including any fields passed
func FatalWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(log.FatalLevel, fields, args...)
}

// ErrorWithFields logs a message at ErrorLevel. Including any fields passed
func ErrorWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(log.ErrorLevel, fields, args...)
}

// WarnWithFields logs a message at WarnLevel. Including any fields passed
func WarningWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(log.WarnLevel, fields, args...)
}

// InfoWithFields logs a message at InfoLevel. Including any fields passed
func InfoWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(log.InfoLevel, fields, args...)
}

// DebugWithFields logs a message at DebugLevel. Including any fields passed
func DebugWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(log.DebugLevel, fields, args...)
}

// TraceWithFields logs a message at TraceLevel. Including any fields passed
func TraceWithFields(fields map[string]interface{}, args ...interface{}) {
	logWithFields(log.TraceLevel, fields, args...)
}

// logWithFields logs a message at the specified level. Including any fields passed
func logWithFields(level log.Level, fields map[string]interface{}, args ...interface{}) {
	if w.logger.GetLevel() >= level {
		entry := w.logger.WithFields(fields)
		if w.reportCaller && level > log.InfoLevel {
			entry.Data["file"] = findCaller()
		}
		entry.Log(level, args...)
	}
}

// findCaller returns the file:line of the caller that started the log call so that the correct file and line number
// are reported in the log message.
func findCaller() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
