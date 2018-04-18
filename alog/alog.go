package alog

import (
	"os"

	"github.com/mushroomsir/logger/pkg"
)

var defaultLogger = pkg.New(os.Stderr, pkg.Options{
	EnableJSON:     true,
	EnableFileLine: true,
})

func checkSugar(v ...interface{}) bool {
	if len(v) >= 2 && len(v)%2 == 0 && v[1] == nil {
		if val, _ := v[0].(string); val == "Error" {
			return false
		}
	}
	return true
}

// SetLevel ...
func SetLevel(level pkg.Level) {
	defaultLogger.SetLevel(level)
}

// Debug ...
func Debug(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Debug(v...)
	}

}

// Info ...
func Info(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Info(v...)
	}
}

// Notice ...
func Notice(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Notice(v...)
	}
}

// Warning ...
func Warning(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Warning(v...)
	}
}

// Check ...
func Check(err interface{}, v ...interface{}) bool {
	return defaultLogger.Check(err, v...)
}

// Err ...
func Err(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Err(v...)
	}
}

// Crit ...
func Crit(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Crit(v...)
	}
}

// Alert ...
func Alert(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Alert(v...)
	}
}

// Emerg ...
func Emerg(v ...interface{}) {
	if checkSugar(v...) {
		defaultLogger.Emerg(v...)
	}
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
