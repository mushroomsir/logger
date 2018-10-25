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

	logger := pkg.New(os.Stderr, pkg.Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	logger.Notice("msg", "content")
	//output: [2017-10-11T02:48:28.598Z] NOTICE {"msg":"content"}
	logger.Err("msg", "content", "code", 500)
	//output: [2017-09-29T05:27:10.639Z] ERR {"code":500,"msg":"content"}
}
