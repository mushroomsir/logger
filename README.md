# logger
[![Build Status](https://img.shields.io/travis/mushroomsir/logger.svg?style=flat-square)](https://travis-ci.org/mushroomsir/logger)
[![Coverage Status](http://img.shields.io/coveralls/mushroomsir/logger.svg?style=flat-square)](https://coveralls.io/github/mushroomsir/logger?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/mushroomsir/logger/blob/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/mushroomsir/logger)

## Installation

```sh
go get github.com/mushroomsir/logger
```

## Usage
```go
package main

import "github.com/mushroomsir/logger"

func main() {
	logger.Debug("xxx")
	//output: [2017-09-29T03:45:11.142Z] DEBUG xxx
	logger.Infof("%v", 1)
	//output: [2017-09-29T03:47:05.436Z] INFO 1
	logger.Warning(logger.Log{"msg": "content", "code": 500})
	//output: [2017-09-29T05:27:10.639Z] WARNING {"code":500,"msg":"content"}
	logger.Warning(map[string]string{"msg": "content"})
	//output: [2017-09-29T05:26:01.638Z] WARNING {"msg":"content"}
	logger.Err(map[string]interface{}{"msg": "content", "code": 500})
	//output: [2017-09-29T05:27:10.639Z] ERR {"code":500,"msg":"content"}
}

```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/logger/blob/master/LICENSE).
