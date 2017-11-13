package main

import (
	"os"

	"github.com/mushroomsir/logger"
)

func main() {
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
