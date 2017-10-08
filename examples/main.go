package main

import (
	"os"

	"github.com/mushroomsir/logger"
)

func main() {
	logger.Debug("xxx")
	//output: [2017-09-29T03:45:11.142Z] DEBUG xxx
	logger.Infof("%v", 1)
	//output: [2017-09-29T03:47:05.436Z] INFO 1
	logger.Warning(logger.Log{"msg": "content", "code": 500})
	//output: [2017-09-29T05:27:10.639Z] WARNING {"code":500,"msg":"content"}

	logger := logger.New(os.Stderr, logger.Options{
		EnableJSON: true,
	})
	logger.Warning(map[string]string{"msg": "content"})
	//output: [2017-09-29T05:26:01.638Z] WARNING {"msg":"content"}
	logger.Err(map[string]interface{}{"msg": "content", "code": 500})
	//output: [2017-09-29T05:27:10.639Z] ERR {"code":500,"msg":"content"}
}
