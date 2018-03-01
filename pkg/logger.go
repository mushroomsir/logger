package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Log ...
type Log map[string]interface{}

// Level represents logging level
// https://tools.ietf.org/html/rfc5424
// https://en.wikipedia.org/wiki/Syslog
type Level uint32

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
	levels = []string{"EMERG", "ALERT", "CRIT", "ERR", "WARNING", "NOTICE", "INFO", "DEBUG"}
)

// Options ...
type Options struct {
	LogFormat      string
	TimeFormat     string
	EnableJSON     bool
	EnableFileLine bool
	Skip           int
}

// New create logger instance
func New(w io.Writer, options ...Options) *Logger {
	logger := &Logger{
		Out:        w,
		enableJSON: true,
		tf:         "2006-01-02T15:04:05.999Z",
		lf:         "[%s] %s %s",
	}
	atomic.StoreUint32(&logger.ulevel, uint32(InfoLevel))
	if len(options) == 0 {
		return logger
	}
	opt := options[0]
	logger.enableJSON = opt.EnableJSON
	if opt.TimeFormat != "" {
		logger.tf = opt.TimeFormat
	}
	if opt.LogFormat != "" {
		logger.lf = opt.LogFormat
	}
	logger.opt = opt
	return logger
}

// Logger ...
type Logger struct {
	Out        io.Writer
	mu         sync.Mutex
	tf, lf     string
	enableJSON bool
	ulevel     uint32
	opt        Options
}

func (a *Logger) checkLogLevel(level Level) bool {
	val := atomic.LoadUint32(&a.ulevel)
	return uint32(level) <= val
}

// SetLevel set the logger's log level
// The default logger level is DebugLevel
func (a *Logger) SetLevel(level Level) {
	ulevel := uint32(7)
	if level > 0 && level <= 7 {
		ulevel = uint32(level)
	}
	atomic.StoreUint32(&a.ulevel, ulevel)
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
		a.Output(time.Now(), DebugLevel, a.format(a.magic(v...)))
	}
}

// Info ...
func (a *Logger) Info(v ...interface{}) {
	if a.checkLogLevel(InfoLevel) {
		a.Output(time.Now(), InfoLevel, a.format(a.magic(v...)))
	}
}

// Notice ...
func (a *Logger) Notice(v ...interface{}) {
	if a.checkLogLevel(NoticeLevel) {
		a.Output(time.Now(), NoticeLevel, a.format(a.magic(v...)))
	}
}

// Warning ...
func (a *Logger) Warning(v ...interface{}) {
	if a.checkLogLevel(WarningLevel) {
		a.Output(time.Now(), WarningLevel, a.format(a.magic(v...)))
	}
}

// Err ...
func (a *Logger) Err(v ...interface{}) {
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now(), ErrLevel, a.format(a.magic(v...)))
	}
}

// Crit ...
func (a *Logger) Crit(v ...interface{}) {
	if a.checkLogLevel(CritiLevel) {
		a.Output(time.Now(), CritiLevel, a.format(a.magic(v...)))
	}
}

// Alert ...
func (a *Logger) Alert(v ...interface{}) {
	if a.checkLogLevel(AlertLevel) {
		a.Output(time.Now(), AlertLevel, a.format(a.magic(v...)))
	}
}

// Emerg ...
func (a *Logger) Emerg(v ...interface{}) {
	if a.checkLogLevel(EmergLevel) {
		a.Output(time.Now(), EmergLevel, a.format(a.magic(v...)))
	}
}
func (a *Logger) magic(v ...interface{}) interface{} {
	if len(v) == 0 {
		return ""
	}
	if len(v) == 1 {
		return v[0]
	}
	m := make(map[string]interface{})
	if len(v)%2 == 0 && a.enableJSON {
		for i, val := range v {
			if i%2 == 0 {
				switch val.(type) {
				case string:
					rVal := v[i+1]
					switch rVal.(type) {
					case error:
						m[val.(string)] = rVal.(error).Error()
					default:
						m[val.(string)] = rVal
					}
				default:
					goto join
				}
			}
		}
		if a.opt.EnableFileLine {
			if a.opt.Skip == 0 {
				a.opt.Skip = 4
			}
			m["FileLine"] = GetCaller(a.opt.Skip)
		}
		return m
	}
join:
	if a.enableJSON {
		m = make(map[string]interface{})
		for i, val := range v {
			m["msg"+strconv.Itoa(i+1)] = val
		}
		return m
	}
	return fmt.Sprint(v...)
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
	}
	switch v.(type) {
	case Log:
		isMarshal = true
	case map[string]interface{}:
		isMarshal = true
	case string:
		return v.(string)
	}
	if isMarshal {
		res, err := json.Marshal(v)
		if err == nil {
			return string(res[:])
		}
	}
	return fmt.Sprint(v)
}

// GetCaller ...
func GetCaller(layer int) string {
	_, file, line, ok := runtime.Caller(layer)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return ""
}

// Stack formats a stack trace of the calling goroutine
func Stack() string {
	buf := make([]byte, 4098)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
