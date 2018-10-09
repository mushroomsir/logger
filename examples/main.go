package main

import (
	"errors"
	"os"

	"github.com/mushroomsir/logger"
	"github.com/mushroomsir/logger/alog"
)

func main() {
	alog.Check(nil)
	alog.Check(errors.New("error"))
	// Simple model
	alog.Info("hello world")
	alog.Info("hello world", " mushroom")
	alog.Infof("hello world %v", "format")

	alog.Info("key", "val")
	alog.Info("Error", nil)
	alog.Info(1, "x")
	alog.Info("Error", errors.New("EOF"))

	// sugar
	logger.Debug("xxx")

	//output: [2017-09-29T03:45:11.142Z] DEBUG xxx
	logger.Infof("%v", 1)
	//output: [2017-09-29T03:47:05.436Z] INFO 1
	logger.Warning("msg", "content", "code", 500)
	//output: [2017-09-29T05:27:10.639Z] WARNING {"code":500,"msg":"content"}

	logger := logger.New(os.Stderr, logger.Options{
		EnableJSON:     true,
		EnableFileLine: true,
	})
	logger.Notice("msg", "content")
	//output: [2017-10-11T02:48:28.598Z] NOTICE {"msg":"content"}
	logger.Err("msg", "content", "code", 500)
	//output: [2017-09-29T05:27:10.639Z] ERR {"code":500,"msg":"content"}
}
