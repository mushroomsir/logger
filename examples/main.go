package main

import (
	"errors"
	"os"

	"github.com/mushroomsir/logger/alog"
	"github.com/mushroomsir/logger/pkg"
)

func main() {
	alog.Info()
	alog.Info(nil)
	alog.IsNil(errors.New("is nil"))
	alog.Info(errors.New("xxx"))
	alog.NotNil(errors.New("error"))
	// Simple model
	alog.Info("hello world")
	alog.Info("hello world", " mushroom")
	alog.Infof("hello world %v", "format")

	alog.Info("key", "val")
	alog.Info("Error", nil)
	alog.Info(1, "x")
	alog.Info("Error", errors.New("EOF"))

	alog.Info(map[string]interface{}{
		"intstr": 1,
	})
	alog.Infof("display message field")

	logger := pkg.New(os.Stderr, pkg.Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	logger.Notice("msg", "content")
	logger.Err("msg", "content", "code", 500)
}
