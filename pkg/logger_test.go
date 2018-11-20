package pkg

import (
	"bytes"
	"errors"
	"math"
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

	require.Contains(buf.String(), `ERR {"error":"x","file":"logger_test.go:51`)
	buf.Reset()

	require.Equal(true, logger.NotNil(errors.New("invalid args"), "Userid", "123456"))
	require.Contains(buf.String(), `"error":"invalid args","file":"logger_test.go:56`)
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
	null := defaultLogger.magic()
	require.Contains(null, `,"message":null}`)
	require.True(defaultLogger.NotNil("nil", "key", "value"))

	null = defaultLogger.magic(1, "X")
	require.Contains(null, `,"message1":1,"message2":"X"}`)
	null = defaultLogger.magic(1, "c", "b")
	require.Contains(null, `,"message1":1,"message2":"c","message3":"b"}`)

	defaultLogger = New(buf, Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	empty := defaultLogger.magic("")
	require.Contains(empty, `,"message1":""}`)
	buf.Reset()

	defaultLogger.Output(time.Now(), DebugLevel, "x\n")
	require.Contains(buf.String(), `] DEBUG x`)

	buf.Reset()
	log := defaultLogger.magic(map[string]interface{}{"x": 1})
	require.Contains(log, `,"x":1`)

	buf.Reset()
	log = defaultLogger.magic("key", "val")
	require.Contains(log, `,"key":"val"}`)

	buf.Reset()
	log = defaultLogger.magic("key", math.NaN())
	require.Equal(log, `{"json-marshal-error":json: unsupported value: NaN}`)

	defaultLogger = New(buf, Options{
		EnableJSON:     false,
		EnableFileLine: true,
	})
	buf.Reset()
	log = defaultLogger.magic("X")
	require.Equal("X", log)
	log = defaultLogger.magic(1)
	require.Equal("1", log)
	log = defaultLogger.magic(1, "X")
	require.Equal("1X", log)

	require.NotEmpty(Stack())

	require.Equal("can not find source file:0", GetCaller(100))

}
