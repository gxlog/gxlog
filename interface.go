package gxlog

import (
	"time"
)

// The Level defines the level type of logs.
type Level int

// All available levels of logs here.
const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Off
)

// LevelCount is the total count of available levels of logs except for Off.
const LevelCount = 6

// A Context is a pair of key-value that is associated with a log Record.
type Context struct {
	Key   string
	Value string
}

// An Auxiliary is a set of extra attributes that are associated with a log Record.
type Auxiliary struct {
	Prefix   string
	Contexts []Context
	Marked   bool
}

// A Record contains all the information of a log.
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

// Formatter is the interface that a formatter of a Logger needs to implement.
// A Formatter must NOT modify the record. In case of asynchrony, a Formatter
// needs to make a new byte slice each time.
//
// Do NOT call methods of the Logger within Format, or it will deadlock.
type Formatter interface {
	Format(record *Record) []byte
}

// Writer is the interface that a writer of a Logger needs to implement.
// A Writer must NOT modify the record.
//
// Do NOT call methods of the Logger within Write, or it will deadlock.
type Writer interface {
	Write(bs []byte, record *Record)
}
