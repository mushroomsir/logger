package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// log ...
type log map[string]interface{}

// Level represents logging level
// https://tools.ietf.org/html/rfc5424
// https://en.wikipedia.org/wiki/Syslog

const (

	// EmergLevel is 0, "Emergency", system is unusable
	EmergLevel uint32 = iota
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
	levels = map[uint32]string{
		0: "EMERG",
		1: "ALERT",
		2: "CRIT",
		3: "ERR",
		4: "WARNING",
		5: "NOTICE",
		6: "INFO",
		7: "DEBUG"}
	message = "message"
	file    = "file"
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
		Out: w,
		tf:  "2006-01-02T15:04:05.999Z",
		lf:  "[%s] %s %s",
	}
	atomic.StoreUint32(&logger.ulevel, InfoLevel)
	if len(options) == 0 {
		return logger
	}
	opt := options[0]
	logger.enableFileLine = opt.EnableFileLine
	logger.enableJSON = opt.EnableJSON
	logger.skip = opt.Skip
	if logger.skip == 0 {
		logger.skip = 3
	}
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
	Out            io.Writer
	mu             sync.Mutex
	tf, lf         string
	enableJSON     bool
	ulevel         uint32
	enableFileLine bool
	skip           int
}

func (a *Logger) checkLogLevel(level uint32) bool {
	val := atomic.LoadUint32(&a.ulevel)
	return level <= val
}

// Level ...
func (a *Logger) Level() uint32 {
	val := atomic.LoadUint32(&a.ulevel)
	return val
}

// SetLevel set the logger's log level
// The default logger level is DebugLevel
func (a *Logger) SetLevel(level uint32) {
	ulevel := InfoLevel
	if level >= 0 && level <= 7 {
		ulevel = level
	}
	atomic.StoreUint32(&a.ulevel, ulevel)
}

// Output ...
func (a *Logger) Output(t time.Time, level uint32, s string) (err error) {
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
func (a *Logger) Debug(kv ...interface{}) {
	if a.checkLogLevel(DebugLevel) {
		a.Output(time.Now(), DebugLevel, a.magic(kv...))
	}
}

// Info ...
func (a *Logger) Info(kv ...interface{}) {
	if a.checkLogLevel(InfoLevel) {
		a.Output(time.Now(), InfoLevel, a.magic(kv...))
	}
}

// Notice ...
func (a *Logger) Notice(kv ...interface{}) {
	if a.checkLogLevel(NoticeLevel) {
		a.Output(time.Now(), NoticeLevel, a.magic(kv...))
	}
}

// Warning ...
func (a *Logger) Warning(kv ...interface{}) {
	if a.checkLogLevel(WarningLevel) {
		a.Output(time.Now(), WarningLevel, a.magic(kv...))
	}
}

// NotNil ...
func (a *Logger) NotNil(err interface{}, kv ...interface{}) bool {
	if IsNil(err) {
		return false
	}
	l := []interface{}{"error", err}
	for _, p := range kv {
		l = append(l, p)
	}
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now(), ErrLevel, a.magic(l...))
	}
	return true
}

// Err ...
func (a *Logger) Err(kv ...interface{}) {
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now(), ErrLevel, a.magic(kv...))
	}
}

// Crit ...
func (a *Logger) Crit(kv ...interface{}) {
	if a.checkLogLevel(CritiLevel) {
		a.Output(time.Now(), CritiLevel, a.magic(kv...))
	}
}

// Alert ...
func (a *Logger) Alert(kv ...interface{}) {
	if a.checkLogLevel(AlertLevel) {
		a.Output(time.Now(), AlertLevel, a.magic(kv...))
	}
}

// Emerg ...
func (a *Logger) Emerg(kv ...interface{}) {
	if a.checkLogLevel(EmergLevel) {
		a.Output(time.Now(), EmergLevel, a.magic(kv...))
	}
}
func (a *Logger) magic(kv ...interface{}) string {
	if !a.enableJSON {
		return fmt.Sprint(kv...)
	}
	m := log{}
	if a.enableFileLine {
		m[file] = GetCaller(a.skip)
	}
	if len(kv) == 0 {
		m[message] = nil
		return a.jsonStr(m)
	}
	if len(kv) == 1 {
		if val, ok := kv[0].(map[string]interface{}); ok {
			for k, v := range val {
				m[k] = v
			}
			return a.jsonStr(m)
		}
	}
	if len(kv)%2 == 0 {
		m1 := log{}
		for i, val := range kv {
			if i%2 == 0 {
				key, ok := val.(string)
				if !ok {
					goto kv
				}
				rVal := kv[i+1]
				switch rVal.(type) {
				case error:
					m1[key] = rVal.(error).Error()
				default:
					m1[key] = rVal
				}
			}
		}
		for k, v := range m1 {
			m[k] = v
		}
		goto json
	}
kv:
	for i, val := range kv {
		m[message+strconv.Itoa(i+1)] = val
	}
json:
	return a.jsonStr(m)
}
func (a *Logger) jsonStr(m log) string {
	res, err := json.Marshal(m)
	if err != nil {
		res = []byte(fmt.Sprintf(`{"json-marshal-error":%v}`, err.Error()))
	}
	return string(res)
}

// Debugf ...
func (a *Logger) Debugf(format string, args ...interface{}) {
	if a.checkLogLevel(DebugLevel) {
		a.Output(time.Now(), DebugLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Infof ...
func (a *Logger) Infof(format string, args ...interface{}) {
	if a.checkLogLevel(InfoLevel) {
		a.Output(time.Now(), InfoLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Noticef ...
func (a *Logger) Noticef(format string, args ...interface{}) {
	if a.checkLogLevel(NoticeLevel) {
		a.Output(time.Now(), NoticeLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Warningf ...
func (a *Logger) Warningf(format string, args ...interface{}) {
	if a.checkLogLevel(WarningLevel) {
		a.Output(time.Now(), WarningLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Errf ...
func (a *Logger) Errf(format string, args ...interface{}) {
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now(), ErrLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Critf ...
func (a *Logger) Critf(format string, args ...interface{}) {
	if a.checkLogLevel(CritiLevel) {
		a.Output(time.Now(), CritiLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Alertf ...
func (a *Logger) Alertf(format string, args ...interface{}) {
	if a.checkLogLevel(AlertLevel) {
		a.Output(time.Now(), AlertLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Emergf ...
func (a *Logger) Emergf(format string, args ...interface{}) {
	if a.checkLogLevel(EmergLevel) {
		a.Output(time.Now(), EmergLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

var workingDir []string

func init() {
	wd, err := os.Getwd()
	if err == nil {
		dir := filepath.ToSlash(wd)
		workingDir = strings.Split(dir, "/")
	}
}

// GetCaller ...
func GetCaller(layer int) string {
	_, file, line, ok := runtime.Caller(layer)
	if !ok {
		file = "can not find source file"
		line = 0
	}
	for _, d := range workingDir {
		if d != "" {
			file = strings.TrimPrefix(file, "/")
			file = strings.TrimPrefix(file, d)
			file = strings.TrimPrefix(file, "/")
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// Stack formats a stack trace of the calling goroutine
func Stack() string {
	buf := make([]byte, 4098)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
