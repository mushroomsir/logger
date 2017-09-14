package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	buf := new(bytes.Buffer)
	logger := New(buf)

	cases := []struct {
		fun   func(v ...interface{})
		funf  func(format string, args ...interface{})
		level string
	}{
		{logger.Debug, logger.Debugf, "DEBUG"},
	}
	for _, c := range cases {
		c.fun("xxx\n")
		assert.True(strings.HasSuffix(buf.String(), c.level+" xxx\n"))
		buf.Reset()
		c.funf("%v %s", "xxx", "1")
		assert.True(strings.HasSuffix(buf.String(), c.level+" xxx 1\n"))
		buf.Reset()
	}
}
