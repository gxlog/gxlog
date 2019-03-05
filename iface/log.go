package iface

import (
	"time"
)

// The Level defines the level type.
type Level int

// All available levels here.
const (
	Trace Level = iota + 1
	Debug
	Info
	Warn
	Error
	Fatal
	Off
)

// LevelCount is the total count of available levels except for Off.
const LevelCount = 6

// A Context is a pair of key-value that is associated with a log.
type Context struct {
	Key   string
	Value string
}

// An Auxiliary is a set of extra attributes that are associated with a log.
type Auxiliary struct {
	Prefix   string
	Contexts []Context
	Marked   bool
}

// A Record contains all information of a log.
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
// needs to make and return a new byte slice each time.
//
// Do NOT call any method of the Logger within Format, or it may deadlock.
type Formatter interface {
	Format(record *Record) []byte
}

// Writer is the interface that a writer of a Logger needs to implement.
// A Writer must NOT modify the bs and record.
//
// Do NOT call any method of the Logger within Write, or it may deadlock.
type Writer interface {
	Write(bs []byte, record *Record)
}
