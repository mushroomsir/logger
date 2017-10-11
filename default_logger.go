package logger

import "os"

var defaultLogger = New(os.Stderr)

// Debug ...
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)

}

// Info ...
func Info(v ...interface{}) {
	defaultLogger.Info(v...)

}

// Notice ...
func Notice(v ...interface{}) {
	defaultLogger.Notice(v...)

}

// Warning ...
func Warning(v ...interface{}) {
	defaultLogger.Warning(v...)

}

// Err ...
func Err(v ...interface{}) {
	defaultLogger.Err(v...)

}

// Crit ...
func Crit(v ...interface{}) {
	defaultLogger.Crit(v...)
}

// Alert ...
func Alert(v ...interface{}) {
	defaultLogger.Alert(v...)
}

// Emerg ...
func Emerg(v ...interface{}) {
	defaultLogger.Emerg(v...)
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
