# logger
[![Build Status](https://img.shields.io/travis/mushroomsir/logger.svg?style=flat-square)](https://travis-ci.org/mushroomsir/logger)
[![Coverage Status](http://img.shields.io/coveralls/mushroomsir/logger.svg?style=flat-square)](https://coveralls.io/github/mushroomsir/logger?branch=master)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/mushroomsir/logger/blob/master/LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/mushroomsir/logger)

## Installation

```sh
go get -u github.com/mushroomsir/logger
```

## Feature

- Easy to use
- Display log ```FileLine```
- Output ```JSON``` format 
- Support ```KV```  syntactic sugar
- Output ```Err``` automatically if ``` err!=nil ```
- Standard log level [Syslog](https://en.wikipedia.org/wiki/Syslog)
- Flexible for custom
- Control output by level

## Usage

#### Easy to use

```go
alog.Info("hello world")
alog.Infof("hello world %v", "format")
// Outout:
[2018-04-12T14:37:55.272Z] INFO hello world
[2018-04-12T14:37:55.272Z] INFO hello world format
```

####  ```KV``` sugar / FileLine / JSON

```go
alog.Info("key", "val")
// Output:
[2018-04-12T14:46:58.088Z] INFO {"FileLine":"D:/go/src/github.com/mushroomsir/logger/examples/main.go:15","Key":"val"}
```

#### Output ```Err``` automatically if ``` err!=nil ```

```go
alog.Info("Error", nil)
// Output:
Does not output anything if ```Err```==nil

alog.Info("Error", errors.New("EOF"))
[2018-04-12T14:51:41.19Z] INFO {"Error":"EOF","FileLine":"D:/go/src/github.com/mushroomsir/logger/examples/main.go:18"}
```

#### Standard log level [Syslog](https://en.wikipedia.org/wiki/Syslog)

```go
alog.Debug()
alog.Info()
alog.Warning()
...
```

#### Flexible for custom

```go
var blog = pkg.New(os.Stderr, pkg.Options{
	EnableJSON:     true,
	EnableFileLine: true,
    TimeFormat: "2006-01-02T15:04:05.999Z",
    LogFormat: "[%s] %s %s",
})
```

#### Control output by level

```go
alog.SetLevel(pkg.InfoLevel)
```

## Licenses

All source code is licensed under the [MIT License](https://github.com/mushroomsir/logger/blob/master/LICENSE).
