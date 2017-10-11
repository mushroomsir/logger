package logger

import (
	"bytes"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(Stack())
	Warning("x", "x")
	Warning()
	assert.Equal(uint32(7), atomic.LoadUint32(&defaultLogger.ulevel))
	defaultLogger.SetLevel(ErrLevel)
	assert.Equal(uint32(3), atomic.LoadUint32(&defaultLogger.ulevel))
	Warning("x") // Don't output
	defaultLogger.SetLevel(DebugLevel)
	assert.Equal(uint32(7), atomic.LoadUint32(&defaultLogger.ulevel))
}
func TestDefaultLogger(t *testing.T) {

	cases := []struct {
		fun   func(v ...interface{})
		funf  func(format string, args ...interface{})
		level string
	}{
		{Debug, Debugf, "DEBUG"},
		{Info, Infof, "INFO"},
		{Notice, Noticef, "NOTICE"},
		{Warning, Warningf, "WARNING"},
		{Err, Errf, "ERR"},
		{Crit, Critf, "CRIT"},
		{Alert, Alertf, "ALERT"},
		{Emerg, Emergf, "EMERG"},
	}

	for _, c := range cases {
		c.fun("hello world\n")
		c.funf("%v %s", "hello world", "1")
	}
}

func TestNoJSON(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	logger := New(buf, Options{
		LogFormat:  "[%s] %s %s",
		TimeFormat: "2006-01-02T15:04:05.999Z",
		EnableJSON: false,
	})

	cases := []struct {
		fun   func(v ...interface{})
		funf  func(format string, args ...interface{})
		level string
	}{
		{logger.Debug, logger.Debugf, "DEBUG"},
		{logger.Info, logger.Infof, "INFO"},
		{logger.Notice, logger.Noticef, "NOTICE"},
		{logger.Warning, logger.Warningf, "WARNING"},
		{logger.Err, logger.Errf, "ERR"},
		{logger.Crit, logger.Critf, "CRIT"},
		{logger.Alert, logger.Alertf, "ALERT"},
		{logger.Emerg, logger.Emergf, "EMERG"},
	}

	for _, c := range cases {
		c.fun("hello world\n")

		assert.True(strings.HasSuffix(buf.String(), "Z] "+c.level+" hello world\n"))
		buf.Reset()

		c.funf("%v %s", "hello world", "1")
		assert.True(strings.HasSuffix(buf.String(), "Z] "+c.level+" hello world 1\n"))
		buf.Reset()

		c.fun(Log{"msg": 1})
		assert.True(strings.HasSuffix(buf.String(), c.level+" {\"msg\":1}\n"))
		buf.Reset()

		c.fun(1)
		assert.True(strings.HasSuffix(buf.String(), "Z] "+c.level+" 1\n"))
		buf.Reset()

		c.fun("x")
		assert.True(strings.HasSuffix(buf.String(), c.level+" x\n"))
		buf.Reset()

		c.fun("name", "vidar", "age", 18)
		assert.True(strings.HasSuffix(buf.String(), c.level+" namevidarage18\n"))
		buf.Reset()

		c.fun(1, "vidar", "age", 18)
		assert.True(strings.HasSuffix(buf.String(), c.level+" 1vidarage18\n"))
		buf.Reset()
	}
}
func TestLogger(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	logger := New(buf, Options{
		LogFormat:  "[%s] %s %s",
		TimeFormat: "2006-01-02T15:04:05.999Z",
		EnableJSON: true,
	})

	cases := []struct {
		fun   func(v ...interface{})
		funf  func(format string, args ...interface{})
		level string
	}{
		{logger.Debug, logger.Debugf, "DEBUG"},
		{logger.Info, logger.Infof, "INFO"},
		{logger.Notice, logger.Noticef, "NOTICE"},
		{logger.Warning, logger.Warningf, "WARNING"},
		{logger.Err, logger.Errf, "ERR"},
		{logger.Crit, logger.Critf, "CRIT"},
		{logger.Alert, logger.Alertf, "ALERT"},
		{logger.Emerg, logger.Emergf, "EMERG"},
	}

	for _, c := range cases {

		c.fun(Log{"msg": 1})
		assert.True(strings.HasSuffix(buf.String(), c.level+" {\"msg\":1}\n"))
		buf.Reset()

		c.fun(map[string]interface{}{"msg": 1})
		assert.True(strings.HasSuffix(buf.String(), c.level+" {\"msg\":1}\n"))
		buf.Reset()

		c.fun(map[string]string{"msg": "1"})
		assert.True(strings.HasSuffix(buf.String(), c.level+" {\"msg\":\"1\"}\n"))
		buf.Reset()

		c.fun(1)
		assert.True(strings.HasSuffix(buf.String(), c.level+" 1\n"))
		buf.Reset()

		c.fun("x")
		assert.True(strings.HasSuffix(buf.String(), c.level+" x\n"))
		buf.Reset()

		c.fun("name", "vidar", "age", 18)
		assert.True(strings.HasSuffix(buf.String(), c.level+" {\"age\":18,\"name\":\"vidar\"}\n"))
		buf.Reset()

		c.fun(1, "vidar", "age", 18)
		assert.True(strings.HasSuffix(buf.String(), c.level+" {\"msg1\":1,\"msg2\":\"vidar\",\"msg3\":\"age\",\"msg4\":18}\n"))
		buf.Reset()
	}
}
