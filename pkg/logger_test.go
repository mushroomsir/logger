package pkg

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type myError struct{}

func (e *myError) Error() string {
	return "x"
}
func GetMyError() interface{} {
	var myerr *myError
	return myerr
}
func GetMyNormalError() error {
	var myerr *myError
	return myerr
}
func ReplaceError(err error) error {
	return err
}
func TestNotNil(t *testing.T) {
	require := require.New(t)
	buf := new(bytes.Buffer)
	logger := New(buf, Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	require.Equal(InfoLevel, logger.Level())
	logger.SetLevel(DebugLevel)
	require.Equal(false, logger.NotNil(nil))
	require.Equal(true, logger.NotNil(errors.New("invalid args")))
	var myerr *myError
	require.True(myerr == nil)

	require.True(ReplaceError(myerr) != nil)
	require.True(GetMyError() != nil)
	require.True(GetMyNormalError() != nil)

	require.Equal(false, logger.NotNil(ReplaceError(myerr)))
	require.Equal(false, logger.NotNil(GetMyError()))
	require.Equal(false, logger.NotNil(GetMyNormalError()))

	require.Equal(true, logger.NotNil(&myError{}))

	require.Contains(buf.String(), `ERR {"error":"x","file":"logger_test.go:50`)
	buf.Reset()

	require.Equal(true, logger.NotNil(errors.New("invalid args"), "Userid", "123456"))
	require.Contains(buf.String(), `"error":"invalid args","file":"logger_test.go:55`)
	require.Contains(buf.String(), `"Userid":"123456"`)
	buf.Reset()

	require.Equal(DebugLevel, logger.Level())
}

func TestDefault(t *testing.T) {
	require := require.New(t)
	buf := new(bytes.Buffer)
	logger := New(buf, Options{
		EnableFileLine: true,
		TimeFormat:     "xxx",
		LogFormat:      "yyy",
	})

	require.Equal("xxx", logger.tf)
	require.Equal("yyy", logger.lf)
	logger = New(buf)
	require.Equal(false, logger.enableFileLine)
}

func TestLevel(t *testing.T) {
	require := require.New(t)

	buf := new(bytes.Buffer)
	defaultLogger := New(buf, Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	defaultLogger.SetLevel(DebugLevel)

	cases := []struct {
		fun   func(v ...interface{})
		funf  func(format string, args ...interface{})
		level string
	}{
		{defaultLogger.Debug, defaultLogger.Debugf, "DEBUG"},
		{defaultLogger.Info, defaultLogger.Infof, "INFO"},
		{defaultLogger.Notice, defaultLogger.Noticef, "NOTICE"},
		{defaultLogger.Warning, defaultLogger.Warningf, "WARNING"},
		{defaultLogger.Err, defaultLogger.Errf, "ERR"},
		{defaultLogger.Crit, defaultLogger.Critf, "CRIT"},
		{defaultLogger.Alert, defaultLogger.Alertf, "ALERT"},
		{defaultLogger.Emerg, defaultLogger.Emergf, "EMERG"},
	}

	for _, c := range cases {
		c.fun("error", "x")
		require.Contains(buf.String(), c.level+` {"error":"x","file":"`)
		buf.Reset()

		c.fun("error", nil, "msg", "xx")
		require.Contains(buf.String(), c.level+` {"error":null,"file":"`)
		buf.Reset()
		c.fun("error", errors.New("xxx"), "msg", "xx")
		require.Contains(buf.String(), c.level+` {"error":"xxx","file":"`)
		buf.Reset()

		c.funf("ccc")
		require.Contains(buf.String(), "")
	}

	require.Equal(false, defaultLogger.NotNil(nil))
	require.Equal(true, defaultLogger.NotNil(errors.New("error")))
	require.Contains(buf.String(), `ERR {"error":"error","file":"logger_test.go:120`)
	buf.Reset()
	require.Equal(true, defaultLogger.NotNil(errors.New("error"), "Userid", "123456"))
	require.Contains(buf.String(), `"error":"error","file":"`)
	require.Contains(buf.String(), `"Userid":"123456"`)
	buf.Reset()
}

func TestCond(t *testing.T) {
	require := require.New(t)

	buf := new(bytes.Buffer)
	defaultLogger := New(buf, Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	defaultLogger.SetLevel(DebugLevel)
	null := defaultLogger.magic().(map[string]interface{})
	require.Equal(2, len(null))
	require.True(defaultLogger.NotNil("nil", "key", "value"))

	null = defaultLogger.magic(1, "X").(map[string]interface{})
	require.Equal(3, len(null))
	null = defaultLogger.magic(1, "c", "b").(map[string]interface{})
	require.Equal(4, len(null))

	defaultLogger = New(buf, Options{
		EnableJSON:     false,
		EnableFileLine: true,
	})
	empty := defaultLogger.magic("")
	require.Equal("", empty)
	buf.Reset()

	defaultLogger.Output(time.Now(), DebugLevel, "x\n")
	require.Contains(buf.String(), `] DEBUG x`)

	log := defaultLogger.format(Log{"x": 1})
	require.Contains(log, "src/testing/testing.go")
	require.Contains(log, `,"x":1`)
	log = defaultLogger.format("X")
	require.Equal("X", log)
	log = defaultLogger.format(1)
	require.Equal("1", log)

	require.NotEmpty(Stack())

	require.Equal("can not find source file:0", GetCaller(100))
}
