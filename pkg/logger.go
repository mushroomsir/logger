package pkg

import (
	"encoding/json"
	"fmt"
	"io"
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
	goID    = "goID"
)

// Options ...
type Options struct {
	LogFormat      string
	TimeFormat     string
	EnableJSON     bool
	EnableFileLine bool
	EnableGoID     bool
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
	logger.enableGoID = opt.EnableGoID
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
	enableGoID     bool
	skip           int
	json           bool
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
// The default logger level is InfoLevel
func (a *Logger) SetLevel(level uint32) *Logger {
	ulevel := InfoLevel
	if level >= 0 && level <= 7 {
		ulevel = level
	}
	atomic.StoreUint32(&a.ulevel, ulevel)
	return a
}

// SetLoggerLevel set the logger's log level
// The default logger level is InfoLevel
func (a *Logger) SetLoggerLevel(level string) *Logger {
	ulevel := ParseLevel(level)
	atomic.StoreUint32(&a.ulevel, ulevel)
	return a
}

// SetJSONLog set the logger writing JSON string log.
func (a *Logger) SetJSONLog() *Logger {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.json = true
	return a
}

// Output ...
func (a *Logger) Output(t time.Time, level uint32, v interface{}) (err error) {
	logObj := format2Log(v)
	if a.json {
		logObj["timestamp"] = t.Format(a.tf)
		logObj["level"] = levels[level]

		str := a.jsonFormat(logObj)

		a.mu.Lock()
		defer a.mu.Unlock()
		_, err = fmt.Fprint(a.Out, str)
		if err == nil {
			a.Out.Write([]byte{'\n'})
		}
	} else {
		str := a.jsonFormat(logObj)

		a.mu.Lock()
		defer a.mu.Unlock()
		if l := len(str); l > 0 && str[l-1] == '\n' {
			str = str[0 : l-1]
		}
		_, err = fmt.Fprintf(a.Out, a.lf, t.UTC().Format(a.tf), levels[level], str)
		if err == nil {
			a.Out.Write([]byte{'\n'})
		}
	}
	return
}

// Debug ...
func (a *Logger) Debug(kv ...interface{}) {
	if a.checkLogLevel(DebugLevel) {
		a.Output(time.Now().UTC(), DebugLevel, a.magic(kv...))
	}
}

// Info ...
func (a *Logger) Info(kv ...interface{}) {
	if a.checkLogLevel(InfoLevel) {
		a.Output(time.Now().UTC(), InfoLevel, a.magic(kv...))
	}
}

// Notice ...
func (a *Logger) Notice(kv ...interface{}) {
	if a.checkLogLevel(NoticeLevel) {
		a.Output(time.Now().UTC(), NoticeLevel, a.magic(kv...))
	}
}

// Warning ...
func (a *Logger) Warning(kv ...interface{}) {
	if a.checkLogLevel(WarningLevel) {
		a.Output(time.Now().UTC(), WarningLevel, a.magic(kv...))
	}
}

// IsNil ...
func (a *Logger) IsNil(err interface{}, kv ...interface{}) bool {
	return !a.NotNil(err, kv)
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
		a.Output(time.Now().UTC(), ErrLevel, a.magic(l...))
	}
	return true
}

// Err ...
func (a *Logger) Err(kv ...interface{}) {
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now().UTC(), ErrLevel, a.magic(kv...))
	}
}

// Crit ...
func (a *Logger) Crit(kv ...interface{}) {
	if a.checkLogLevel(CritiLevel) {
		a.Output(time.Now().UTC(), CritiLevel, a.magic(kv...))
	}
}

// Alert ...
func (a *Logger) Alert(kv ...interface{}) {
	if a.checkLogLevel(AlertLevel) {
		a.Output(time.Now().UTC(), AlertLevel, a.magic(kv...))
	}
}

// Emerg ...
func (a *Logger) Emerg(kv ...interface{}) {
	if a.checkLogLevel(EmergLevel) {
		a.Output(time.Now().UTC(), EmergLevel, a.magic(kv...))
	}
}

func (a *Logger) magic(kv ...interface{}) interface{} {
	if !a.enableJSON {
		return fmt.Sprint(kv...)
	}
	m := log{}
	if a.enableFileLine {
		m[file] = GetCaller(a.skip)
	}
	if len(kv) == 0 {
		m[message] = nil
		return m
	}
	if len(kv) == 1 {
		if val, ok := kv[0].(map[string]interface{}); ok {
			for k, v := range val {
				m[k] = v
			}
			return m
		}
	}
	if len(kv)%2 == 0 {
		m1 := log{}
		for i, val := range kv {
			if i%2 == 0 {
				key, ok := val.(string)
				if !ok {
					goto kvBlock
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
		goto jsonBlock
	}
kvBlock:
	for i, val := range kv {
		key := message + strconv.Itoa(i+1)
		switch val.(type) {
		case error:
			m[key] = val.(error).Error()
		default:
			m[key] = val
		}
	}
jsonBlock:
	return m
}
func (a *Logger) jsonFormat(m log) string {
	if a.enableGoID {
		m[goID] = GoroutineID()
	}
	res, err := json.Marshal(m)
	if err != nil {
		res = []byte(fmt.Sprintf(`{"json-marshal-error":%v}`, err.Error()))
	}
	return string(res)
}

// Debugf ...
func (a *Logger) Debugf(format string, args ...interface{}) {
	if a.checkLogLevel(DebugLevel) {
		a.Output(time.Now().UTC(), DebugLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Infof ...
func (a *Logger) Infof(format string, args ...interface{}) {
	if a.checkLogLevel(InfoLevel) {
		a.Output(time.Now().UTC(), InfoLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Noticef ...
func (a *Logger) Noticef(format string, args ...interface{}) {
	if a.checkLogLevel(NoticeLevel) {
		a.Output(time.Now().UTC(), NoticeLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Warningf ...
func (a *Logger) Warningf(format string, args ...interface{}) {
	if a.checkLogLevel(WarningLevel) {
		a.Output(time.Now().UTC(), WarningLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Errf ...
func (a *Logger) Errf(format string, args ...interface{}) {
	if a.checkLogLevel(ErrLevel) {
		a.Output(time.Now().UTC(), ErrLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Critf ...
func (a *Logger) Critf(format string, args ...interface{}) {
	if a.checkLogLevel(CritiLevel) {
		a.Output(time.Now().UTC(), CritiLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Alertf ...
func (a *Logger) Alertf(format string, args ...interface{}) {
	if a.checkLogLevel(AlertLevel) {
		a.Output(time.Now().UTC(), AlertLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Emergf ...
func (a *Logger) Emergf(format string, args ...interface{}) {
	if a.checkLogLevel(EmergLevel) {
		a.Output(time.Now().UTC(), EmergLevel, a.magic(message, fmt.Sprintf(format, args...)))
	}
}

// Panicf ...
func (a *Logger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	a.Output(time.Now().UTC(), EmergLevel, a.magic(message, msg))
	panic(msg)
}

// GetCaller ...
func GetCaller(layer int) string {
	_, file, line, ok := runtime.Caller(layer)
	if !ok {
		file = "can not find source file"
		line = 0
	}
	files := strings.Split(file, "/")
	if len(files) > 3 {
		filesLen := len(files)
		file = files[filesLen-3] + "/" + files[filesLen-2] + "/" + files[filesLen-1]
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// Stack formats a stack trace of the calling goroutine
func Stack() string {
	buf := make([]byte, 4098)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func format2Log(i interface{}) log {
	switch v := i.(type) {
	case log:
		return v
	case map[string]interface{}:
		return log(v)
	default:
		return log{"message": fmt.Sprint(i)}
	}
}

// ParseLevel takes a string level and returns the logging level constant.
func ParseLevel(level string) uint32 {
	switch strings.ToLower(level) {
	case "emergency", "emerg":
		return EmergLevel
	case "alert":
		return AlertLevel
	case "critical", "crit", "criti":
		return CritiLevel
	case "error", "err":
		return ErrLevel
	case "warning", "warn":
		return WarningLevel
	case "notice":
		return NoticeLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	}
	return InfoLevel
}
