package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(Stack())
}
func TestDefaultLogger(t *testing.T) {

	cases := []struct {
		fun   func(v interface{})
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
		Level:      DebugLevel,
	})

	cases := []struct {
		fun   func(v interface{})
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
		fun   func(v interface{})
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
	}

}
