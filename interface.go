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
	Debugf(fmt string, args ...interface{})
	Info(args ...interface{})
	Infof(fmt string, args ...interface{})
	Warn(args ...interface{})
	Warnf(fmt string, args ...interface{})
	Error(args ...interface{})
	Errorf(fmt string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(fmt string, args ...interface{})
}

type Record struct {
	Time  time.Time
	Level LogLevel
	Func  string
	Line  int
	Path  string
	File  string
	Msg   []byte
}

type Formatter interface {
	Format(*Record) []byte
}

type Writer interface {
	io.Writer
}
