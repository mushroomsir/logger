package pkg

import (
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

type myError struct{}

func (e *myError) Error() string {
	return ""
}
func TestCheck(t *testing.T) {
	require := require.New(t)
	buf := new(bytes.Buffer)
	logger := New(buf, Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	logger.SetLevel(DebugLevel)

	require.Equal(false, logger.Check(nil))
	require.Equal(true, logger.Check(errors.New("error")))
	require.Equal(true, logger.Check(&myError{}))

	log.Println(buf.String())
	require.Contains(buf.String(), `ERR {"Error":"error","FileLine":"`)
	buf.Reset()

	require.Equal(true, logger.Check(errors.New("error"), "Userid", "123456"))
	log.Println(buf.String())
	require.Contains(buf.String(), `ERR {"Error":"error","FileLine":"`)
	require.Contains(buf.String(), `,"Userid":"123456"`)
	buf.Reset()
}
