package gxlog

import (
	"time"
)

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOff
	LevelCount = LevelOff
)

type Context struct {
	Key   string
	Value string
}

type Auxiliary struct {
	Prefix   string
	Contexts []Context
	Marked   bool
}

type Record struct {
	Time  time.Time
	Level Level
	File  string
	Line  int
	Pkg   string
	Func  string
	Msg   string
	Aux   Auxiliary
}

type Formatter interface {
	Format(record *Record) []byte
}

type Writer interface {
	Write(bs []byte, record *Record)
}
