package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
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
	defaultLevel Level = iota
	// EmergLevel is 0, "Emergency", system is unusable
	EmergLevel
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
	levels = []string{"EMERG", "ALERT", "CRIT", "ERR", "WARNING", "NOTICE", "INFO", "DEBUG"}
)

// Options ...
type Options struct {
	LogFormat  string
	TimeFormat string
	Level      Level
	EnableJSON bool
}

// New create logger instance
func New(w io.Writer, options ...Options) *Logger {
	logger := &Logger{
		Out:        w,
		level:      DebugLevel,
		enableJSON: true,
		tf:         "2006-01-02T15:04:05.999Z",
		lf:         "[%s] %s %s",
	}
	if len(options) == 0 {
		return logger
	}
	opt := options[0]
	if opt.Level != defaultLevel {
		logger.level = opt.Level
	}
	logger.enableJSON = opt.EnableJSON
	if opt.TimeFormat != "" {
		logger.tf = opt.TimeFormat
	}
	if opt.LogFormat != "" {
		logger.lf = opt.LogFormat
	}
	return logger
}

// Logger ...
type Logger struct {
	Out        io.Writer
	mu         sync.Mutex
	tf, lf     string
	level      Level
	enableJSON bool
}

func (a *Logger) checkLogLevel(level Level) bool {
	return level <= a.level
}

// Output ...
func (a *Logger) Output(t time.Time, level Level, s string) (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if l := len(s); l > 0 && s[l-1] == '\n' {
		s = s[0 : l-1]
	}
	_, err = fmt.Fprintf(a.Out, a.lf, t.UTC().Format(a.tf), levels[level-1], s)
	if err == nil {
		a.Out.Write([]byte{'\n'})
	}
	return
}

// Debug ...
func (a *Logger) Debug(v interface{}) {
	if a.checkLogLevel(DebugLevel) {
		a.Output(time.Now(), DebugLevel, a.format(v))
	}
}

// Info ...
func (a *Logger) Info(v interface{}) {
	if a.checkLogLevel(InfoLevel) {
		a.Output(time.Now(), InfoLevel, a.format(v))
	}
}

// Notice ...
func (a *Logger) Notice(v interface{}) {
	if a.checkLogLevel(NoticeLevel) {
		a.Output(time.Now(), NoticeLevel, a.format(v))
	}
}

// Warning ...
func (a *Logger) Warning(v interface{}) {
	if a.checkLogLevel(WarningLevel) {
		a.Output(time.Now(), WarningLevel, a.format(v))
	}
}

// Err ...
func (a *Logger) Err(v interface{}) {
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now(), ErrLevel, a.format(v))
	}
}

// Crit ...
func (a *Logger) Crit(v interface{}) {
	if a.checkLogLevel(CritiLevel) {
		a.Output(time.Now(), CritiLevel, a.format(v))
	}
}

// Alert ...
func (a *Logger) Alert(v interface{}) {
	if a.checkLogLevel(AlertLevel) {
		a.Output(time.Now(), AlertLevel, a.format(v))
	}
}

// Emerg ...
func (a *Logger) Emerg(v interface{}) {
	if a.checkLogLevel(EmergLevel) {
		a.Output(time.Now(), EmergLevel, a.format(v))
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

func (a *Logger) format(v interface{}) string {
	var isMarshal bool
	if a.enableJSON {
		isMarshal = true
	} else {
		switch v.(type) {
		case Log:
			isMarshal = true
		case string:
			return v.(string)
		}
	}
	if isMarshal {
		res, err := json.Marshal(v)
		if err == nil {
			return string(res)
		}
	}
	return fmt.Sprint(v)
}

// Stack formats a stack trace of the calling goroutine
func Stack() string {
	buf := make([]byte, 4098)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
