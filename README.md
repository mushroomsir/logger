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
- Improve efficiency
- Output ```Err``` automatically if err not nil
- Standard log level [Syslog](https://en.wikipedia.org/wiki/Syslog)
- Flexible for custom
- Control output by level

## Usage

#### Easy to use

```go
alog.Info("hello world")
alog.Infof("hello world %v", "format")
// Outout:
[2018-10-13T03:05:28.476Z] INFO {"file":"examples/main.go:16","message1":"hello world"}
[2018-10-13T03:05:28.477Z] INFO {"file":"examples/main.go:18","message":"hello world format"}
```

####  ```KV``` sugar / FileLine / JSON

```go
alog.Info("key", "val")
// Output:
[2018-04-12T14:46:58.088Z] INFO {"file":"main.go:15","Key":"val"}
```
#### Improve efficiency
##### Return ```true``` value and output ```Err``` log if ``` err!=nil ```
```go
err := errors.New("x")
if alog.NotNil(err) {
    return err
}
// Output:
[2018-04-18T00:34:19.946Z] ERR {"error":"x","file":"main.go:13"}
```
##### Does not output anything if Err==nil and continue code execution
```go
var err error
if alog.NotNil(err) {
    return err
}
// continue code execution
```

#### Standard log level [Syslog](https://en.wikipedia.org/wiki/Syslog)

```go
alog.Debug()
alog.Debugf()
alog.Info()
alog.Infof()
alog.Warning()
alog.Warningf()
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
