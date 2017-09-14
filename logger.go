package logger

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/teambition/gear"
)

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

var levels = []string{"EMERG", "ALERT", "CRIT", "ERR", "WARNING", "NOTICE", "INFO", "DEBUG"}

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
		panic(gear.Err.WithMsg("invalid logger level"))
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
func (a *Logger) Debug(v ...interface{}) {
	if a.checkLogLevel(DebugLevel) {
		a.Output(time.Now(), DebugLevel, fmt.Sprint(v...))
	}
}

// Debugf ...
func (a *Logger) Debugf(format string, args ...interface{}) {
	a.Debug(fmt.Sprintf(format, args...))
}
