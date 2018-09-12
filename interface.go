package gxlog

import (
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

type Record struct {
	Time     time.Time
	Level    LogLevel
	Pathname string
	Line     int
	Func     string
	Msg      string
}

type Formatter interface {
	Format(*Record) []byte
}

type Writer interface {
	Write(bs []byte, record *Record)
}
