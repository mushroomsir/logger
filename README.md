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
	// sugar
	logger.Debug("xxx")
	//output: [2017-09-29T03:45:11.142Z] DEBUG xxx
	logger.Infof("%v", 1)
	//output: [2017-09-29T03:47:05.436Z] INFO 1
	logger.Warning("msg", "content", "code", 500)
	//output: [2017-09-29T05:27:10.639Z] WARNING {"code":500,"msg":"content"}

	logger := logger.New(os.Stderr, logger.Options{
		EnableJSON: true,
	})
	logger.Notice("msg", "content")
	//output: [2017-10-11T02:48:28.598Z] NOTICE {"msg":"content"}
	logger.Err("msg", "content", "code", 500)
	//output: [2017-09-29T05:27:10.639Z] ERR {"code":500,"msg":"content"}
}



```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/logger/blob/master/LICENSE).
