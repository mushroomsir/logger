package alog

import (
	"bytes"
	"errors"
	"testing"

	"github.com/mushroomsir/logger/pkg"
	"github.com/stretchr/testify/require"
)

func TestAlog(t *testing.T) {
	require := require.New(t)

	buf := new(bytes.Buffer)
	defaultLogger = pkg.New(buf, pkg.Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	defaultLogger.SetLevel(pkg.DebugLevel)

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
		require.False(checkSugar("Error", nil))
		c.fun("Error", nil)
		require.Equal("", buf.String())
		buf.Reset()

		require.False(checkSugar("Error", nil, "msg", "xx"))
		c.fun("Error", nil, "msg", "xx")
		require.Equal("", buf.String())
		buf.Reset()

		require.True(checkSugar("error", nil, "msg", "xx"))
		c.fun("error", nil, "msg", "xx")
		require.NotEqual("", buf.String())
		buf.Reset()
	}

	require.Equal(false, Check(nil))
	require.Equal(true, Check(errors.New("error")))
	require.Contains(buf.String(), `ERR {"Error":"error","FileLine":"`)
	buf.Reset()
	require.Equal(true, Check(errors.New("error"), "Userid", "123456"))
	require.Contains(buf.String(), `ERR {"Error":"error","FileLine":"`)
	require.Contains(buf.String(), `,"Userid":"123456"`)
	buf.Reset()
}
