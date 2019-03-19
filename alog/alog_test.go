package alog

import (
	"bytes"
	"errors"
	"testing"

	"github.com/mushroomsir/logger/pkg"
	"github.com/stretchr/testify/require"
)

func TestLevel(t *testing.T) {
	require := require.New(t)
	SetLevel(6)
	require.Equal(uint32(6), Level())

	SetLevel(9)
	require.Equal(uint32(6), Level())
}
func TestAlog(t *testing.T) {
	require := require.New(t)

	buf := new(bytes.Buffer)
	defaultLogger = pkg.New(buf, pkg.Options{
		EnableJSON:     true,
		EnableFileLine: true,
		Skip:           4,
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

	require.Equal(false, NotNil(nil))
	require.Equal(true, NotNil(errors.New("error")))
	require.Contains(buf.String(), `ERR {"error":"error","file":"logger/alog/alog_test.go:63`)
	buf.Reset()
	require.Equal(true, NotNil(errors.New("error"), "Userid", "123456"))
	require.Contains(buf.String(), `"error":"error","file":"`)
	require.Contains(buf.String(), `"Userid":"123456"`)
	buf.Reset()

	require.Equal(false, Check(nil))
	require.Equal(true, Check(errors.New("error")))
	require.Contains(buf.String(), `ERR {"error":"error","file":"logger/alog/alog_test.go:72`)
	buf.Reset()
	require.Equal(true, Check(errors.New("error"), "Userid", "123456"))
	require.Contains(buf.String(), `"error":"error","file":"`)
	require.Contains(buf.String(), `"Userid":"123456"`)
	buf.Reset()
}
