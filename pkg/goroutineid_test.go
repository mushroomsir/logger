package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoroutineID(t *testing.T) {
	require := require.New(t)
	require.True(GoroutineID() > 0)
}
