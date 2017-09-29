package logger

import "os"

var defaultLogger = New(os.Stderr)

// Debug ...
func Debug(v interface{}) {
	defaultLogger.Debug(v)

}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Info ...
func Info(v interface{}) {
	defaultLogger.Info(v)

}

// Infof ...
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Notice ...
func Notice(v interface{}) {
	defaultLogger.Notice(v)

}

// Noticef ...
func Noticef(format string, args ...interface{}) {
	defaultLogger.Noticef(format, args...)
}

// Warning ...
func Warning(v interface{}) {
	defaultLogger.Warning(v)

}

// Warningf ...
func Warningf(format string, args ...interface{}) {
	defaultLogger.Warningf(format, args...)
}

// Err ...
func Err(v interface{}) {
	defaultLogger.Err(v)

}

// Errf ...
func Errf(format string, args ...interface{}) {
	defaultLogger.Errf(format, args...)
}

// Crit ...
func Crit(v interface{}) {
	defaultLogger.Crit(v)

}

// Critf ...
func Critf(format string, args ...interface{}) {
	defaultLogger.Critf(format, args...)
}

// Alert ...
func Alert(v interface{}) {
	defaultLogger.Alert(v)

}

// Alertf ...
func Alertf(format string, args ...interface{}) {
	defaultLogger.Alertf(format, args...)
}

// Emerg ...
func Emerg(v interface{}) {
	defaultLogger.Emerg(v)

}

// Emergf ...
func Emergf(format string, args ...interface{}) {
	defaultLogger.Emergf(format, args...)
}
