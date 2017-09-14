package logger

import "os"

var defaultLogger = New(os.Stderr)

// Debug ...
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)

}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}
