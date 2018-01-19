package logger

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {

	logger := New(os.Stderr)

	logger.Warning("Time", time.Now())
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	logger.Warning("Duration", time.Since(start))
	logger.Warning("Seconds", time.Since(start).Seconds())
	logger.Warning("String", time.Since(start).String())
}
