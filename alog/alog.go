package alog

import (
	"os"

	"github.com/mushroomsir/logger/pkg"
)

var defaultLogger = pkg.New(os.Stderr, pkg.Options{
	EnableFileLine: true,
	EnableJSON:     true,
	Skip:           4,
	EnableGoID:     true,
})

// SetLevel ...
func SetLevel(level uint32) {
	defaultLogger.SetLevel(level)
}

// Level ...
func Level() uint32 {
	return defaultLogger.Level()
}

// Debug ...
func Debug(kv ...interface{}) {
	defaultLogger.Debug(kv...)

}

// Info ...
func Info(kv ...interface{}) {
	defaultLogger.Info(kv...)
}

// Notice ...
func Notice(kv ...interface{}) {
	defaultLogger.Notice(kv...)
}

// Warning ...
func Warning(kv ...interface{}) {
	defaultLogger.Warning(kv...)
}

// Check was deprecated please use NotNil
func Check(err interface{}, kv ...interface{}) bool {
	if pkg.IsNil(err) {
		return false
	}
	l := []interface{}{"error", err}
	for _, p := range kv {
		l = append(l, p)
	}
	defaultLogger.Err(l...)
	return true
}

// IsNil ...
func IsNil(err interface{}, kv ...interface{}) bool {
	if pkg.IsNil(err) {
		return true
	}
	l := []interface{}{"error", err}
	for _, p := range kv {
		l = append(l, p)
	}
	defaultLogger.Err(l...)
	return false
}

// NotNil ...
func NotNil(err interface{}, kv ...interface{}) bool {
	if pkg.IsNil(err) {
		return false
	}
	l := []interface{}{"error", err}
	for _, p := range kv {
		l = append(l, p)
	}
	defaultLogger.Err(l...)
	return true
}

// Err ...
func Err(kv ...interface{}) {
	defaultLogger.Err(kv...)
}

// Crit ...
func Crit(kv ...interface{}) {
	defaultLogger.Crit(kv...)
}

// Alert ...
func Alert(kv ...interface{}) {
	defaultLogger.Alert(kv...)
}

// Emerg ...
func Emerg(kv ...interface{}) {
	defaultLogger.Emerg(kv...)
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Noticef ...
func Noticef(format string, args ...interface{}) {
	defaultLogger.Noticef(format, args...)
}

// Warningf ...
func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

// Errf ...
func Errf(format string, args ...interface{}) {
	defaultLogger.Errf(format, args...)
}

// Critf ...
func Critf(format string, args ...interface{}) {
	defaultLogger.Critf(format, args...)
}

// Alertf ...
func Alertf(format string, args ...interface{}) {
	defaultLogger.Alertf(format, args...)
}

// Emergf ...
func Emergf(format string, args ...interface{}) {
	defaultLogger.Emergf(format, args...)
}
