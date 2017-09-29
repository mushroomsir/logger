package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

// Log ...
type Log map[string]interface{}

// Level represents logging level
// https://tools.ietf.org/html/rfc5424
// https://en.wikipedia.org/wiki/Syslog
type Level uint8

const (
	// EmergLevel is 0, "Emergency", system is unusable
	EmergLevel Level = iota
	// AlertLevel is 1, "Alert", action must be taken immediately
	AlertLevel
	// CritiLevel is 2, "Critical", critical conditions
	CritiLevel
	// ErrLevel is 3, "Error", error conditions
	ErrLevel
	// WarningLevel is 4, "Warning", warning conditions
	WarningLevel
	// NoticeLevel is 5, "Notice", normal but significant condition
	NoticeLevel
	// InfoLevel is 6, "Informational", informational messages
	InfoLevel
	// DebugLevel is 7, "Debug", debug-level messages
	DebugLevel
)

var (
	levels          = []string{"EMERG", "ALERT", "CRIT", "ERR", "WARNING", "NOTICE", "INFO", "DEBUG"}
	ErrInvalidLevel = errors.New("invalid logger level")
)

// New create logger instance
func New(w io.Writer) *Logger {
	logger := &Logger{Out: w}
	logger.SetLevel(DebugLevel)
	logger.SetTimeFormat("2006-01-02T15:04:05.999Z")
	logger.SetLogFormat("[%s] %s %s")
	return logger
}

// Logger ...
type Logger struct {
	Out    io.Writer
	mu     sync.Mutex
	tf, lf string
	level  Level
}

func (a *Logger) checkLogLevel(level Level) bool {
	return level <= a.level
}

// SetLevel set the logger's log level
// The default logger level is DebugLevel
func (a *Logger) SetLevel(level Level) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if level > DebugLevel {
		panic(ErrInvalidLevel)
	}
	a.level = level
}

// SetTimeFormat set the logger timestamp format
// The default logger timestamp format is "2006-01-02T15:04:05.999Z"(JavaScript ISO date string)
func (a *Logger) SetTimeFormat(timeFormat string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tf = timeFormat
}

// SetLogFormat set the logger log format
// it should accept 3 string values: timestamp, log level and log message
// The default logger log format is "[%s] %s %s": "[time] logLevel message"
func (a *Logger) SetLogFormat(logFormat string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lf = logFormat
}

// Output ...
func (a *Logger) Output(t time.Time, level Level, s string) (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if l := len(s); l > 0 && s[l-1] == '\n' {
		s = s[0 : l-1]
	}
	_, err = fmt.Fprintf(a.Out, a.lf, t.UTC().Format(a.tf), levels[level], s)
	if err == nil {
		a.Out.Write([]byte{'\n'})
	}
	return
}

// Debug ...
func (a *Logger) Debug(v interface{}) {
	if a.checkLogLevel(DebugLevel) {
		a.Output(time.Now(), DebugLevel, format(v))
	}
}

// Info ...
func (a *Logger) Info(v interface{}) {
	if a.checkLogLevel(InfoLevel) {
		a.Output(time.Now(), InfoLevel, format(v))
	}
}

// Notice ...
func (a *Logger) Notice(v interface{}) {
	if a.checkLogLevel(NoticeLevel) {
		a.Output(time.Now(), NoticeLevel, format(v))
	}
}

// Warning ...
func (a *Logger) Warning(v interface{}) {
	if a.checkLogLevel(WarningLevel) {
		a.Output(time.Now(), WarningLevel, format(v))
	}
}

// Err ...
func (a *Logger) Err(v interface{}) {
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now(), ErrLevel, format(v))
	}
}

// Crit ...
func (a *Logger) Crit(v interface{}) {
	if a.checkLogLevel(CritiLevel) {
		a.Output(time.Now(), CritiLevel, format(v))
	}
}

// Alert ...
func (a *Logger) Alert(v interface{}) {
	if a.checkLogLevel(AlertLevel) {
		a.Output(time.Now(), AlertLevel, format(v))
	}
}

// Emerg ...
func (a *Logger) Emerg(v interface{}) {
	if a.checkLogLevel(EmergLevel) {
		a.Output(time.Now(), EmergLevel, format(v))
	}
}

// Debugf ...
func (a *Logger) Debugf(format string, args ...interface{}) {
	a.Debug(fmt.Sprintf(format, args...))
}

// Infof ...
func (a *Logger) Infof(format string, args ...interface{}) {
	a.Info(fmt.Sprintf(format, args...))
}

// Noticef ...
func (a *Logger) Noticef(format string, args ...interface{}) {
	a.Notice(fmt.Sprintf(format, args...))
}

// Warningf ...
func (a *Logger) Warningf(format string, args ...interface{}) {
	a.Warning(fmt.Sprintf(format, args...))
}

// Errf ...
func (a *Logger) Errf(format string, args ...interface{}) {
	a.Err(fmt.Sprintf(format, args...))
}

// Critf ...
func (a *Logger) Critf(format string, args ...interface{}) {
	a.Crit(fmt.Sprintf(format, args...))
}

// Alertf ...
func (a *Logger) Alertf(format string, args ...interface{}) {
	a.Alert(fmt.Sprintf(format, args...))
}

// Emergf ...
func (a *Logger) Emergf(format string, args ...interface{}) {
	a.Emerg(fmt.Sprintf(format, args...))
}

func format(v interface{}) string {
	var isMarshal bool
	switch v.(type) {
	case map[string]string:
		isMarshal = true
	case map[string]interface{}:
		isMarshal = true
	case Log:
		isMarshal = true
	case string:
		return v.(string)
	}
	if isMarshal {
		res, err := json.Marshal(v)
		if err == nil {
			return string(res)
		}
	}
	return fmt.Sprintln(v)
}
