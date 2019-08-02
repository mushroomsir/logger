package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"testing"

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

	require.Contains(buf.String(), `ERR {"error":"x","file":"logger/pkg/logger_test.go:51`)
	buf.Reset()

	require.Equal(true, logger.NotNil(errors.New("invalid args"), "Userid", "123456"))
	require.Contains(buf.String(), `"error":"invalid args","file":"logger/pkg/logger_test.go:56`)
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
	require.Contains(buf.String(), `ERR {"error":"error","file":"logger/pkg/logger_test.go:120`)
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
	logObj := format2Log(null)
	require.Equal(nil, logObj["message"])
	require.True(defaultLogger.NotNil("nil", "key", "value"))

	null = defaultLogger.magic(1, "X")
	logObj = format2Log(null)
	require.Equal(1, logObj["message1"])
	require.Equal("X", logObj["message2"])

	null = defaultLogger.magic(1, "c", "b")
	logObj = format2Log(null)
	require.Equal(1, logObj["message1"])
	require.Equal("c", logObj["message2"])
	require.Equal("b", logObj["message3"])

	defaultLogger = New(buf, Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	empty := defaultLogger.magic("")
	logObj = format2Log(empty)
	require.Equal("", logObj["message1"])
	buf.Reset()

	buf.Reset()
	log := defaultLogger.magic(map[string]interface{}{"x": 1})
	logObj = format2Log(log)
	require.Equal(1, logObj["x"])

	buf.Reset()
	log = defaultLogger.magic("key", "val")
	logObj = format2Log(log)
	require.Equal("val", logObj["key"])

	buf.Reset()
	log = defaultLogger.magic("key", math.NaN())
	logObj = format2Log(log)
	require.NotNil(logObj["key"])

	defaultLogger = New(buf, Options{
		EnableJSON:     false,
		EnableFileLine: true,
	})
	buf.Reset()
	log = defaultLogger.magic("X")
	logObj = format2Log(log)
	require.Equal("X", logObj["message"])

	log = defaultLogger.magic(1)
	logObj = format2Log(log)
	require.Equal("1", logObj["message"])

	log = defaultLogger.magic(1, "X")
	logObj = format2Log(log)
	require.Equal("1X", logObj["message"])

	require.NotEmpty(Stack())

	require.Equal("can not find source file:0", GetCaller(100))

}

func TestJsonLog(t *testing.T) {
	require := require.New(t)
	buf := new(bytes.Buffer)
	defaultLogger := New(buf, Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	defaultLogger.SetJSONLog()
	defaultLogger.Info("a", "1")

	v := map[string]string{}
	err := json.Unmarshal(buf.Bytes(), &v)
	require.Nil(err)
	require.Equal("1", v["a"])
	require.Equal("INFO", v["level"])
	require.NotEmpty(v["timestamp"])
	require.NotEmpty(v["file"])
}
