package gxlog

import (
	"io"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
)

type LeveledLogger interface {
	Debug(args ...interface{})
	Debugf(fmtstr string, args ...interface{})
	Info(args ...interface{})
	Infof(fmtstr string, args ...interface{})
	Warn(args ...interface{})
	Warnf(fmtstr string, args ...interface{})
	Error(args ...interface{})
	Errorf(fmtstr string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(fmtstr string, args ...interface{})
}

type Record struct {
	Time     time.Time
	Level    LogLevel
	Pathname string
	Line     int
	Func     string
	Msg      []byte
}

type Formatter interface {
	Format(*Record) []byte
}

type Writer interface {
	io.Writer
}
