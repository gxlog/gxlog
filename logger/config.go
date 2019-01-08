package logger

import (
	"github.com/gxlog/gxlog/iface"
)

// The Flag defines the flag type of Logger.
type Flag int

// All available flags here.
const (
	Prefix Flag = 0x01 << iota
	StaticContext
	DynamicContext
	Mark
	LimitByCount
	LimitByTime
)

// The Filter type defines a function type which is used to filter logs.
//
// Do NOT call any method of the Logger within a filter, or it may deadlock.
type Filter func(*iface.Record) bool

// A Config is used to configure a Logger.
type Config struct {
	// Level is the level of Logger, logs with a lower level will be omitted.
	// If it is not specified, Trace is used. Otherwise, its value MUST be
	// between Trace and Off inclusive.
	Level iface.Level
	// TrackLevel is the auto backtracking level of Logger.
	// If the level of a emitted log is NOT lower than the TrackLevel, the stack
	// of the current goroutine will be output with the log.
	// If it is not specified, Fatal is used. Otherwise, its value MUST be
	// between Trace and Off inclusive.
	TrackLevel iface.Level
	// ExitLevel is the auto exiting level of Logger.
	// If the level of a emitted log is NOT lower than the ExitLevel, the Logger
	// will call os.Exit after outputting the log.
	// If it is not specified, Off is used. Otherwise, its value MUST be
	// between Trace and Off inclusive.
	ExitLevel iface.Level
	// TimeLevel is the level of a emitted log when Logger.Time or Logger.Timef
	// is called.
	// If it is not specified, Trace is used. Otherwise, its value MUST be
	// between Trace and Fatal inclusive.
	TimeLevel iface.Level
	// PanicLevel is the level of a emitted log when Logger.Panic or
	// Logger.Panicf is called.
	// If it is not specified, Fatal is used. Otherwise, its value MUST be
	// between Trace and Fatal inclusive.
	PanicLevel iface.Level
	// Filter is the log filter of Logger. If it is not nil, it will be called
	// when a log emits. And if it returns false, the log will be omitted.
	Filter Filter
	// Disabled is a set of flags. If a flag is set, the corresponding feature
	// of Logger will be disabled.
	Disabled Flag
}

func (config *Config) setDefaults() {
	if config.Level == 0 {
		config.Level = iface.Trace
	}
	if config.TrackLevel == 0 {
		config.TrackLevel = iface.Fatal
	}
	if config.ExitLevel == 0 {
		config.ExitLevel = iface.Off
	}
	if config.TimeLevel == 0 {
		config.TimeLevel = iface.Trace
	}
	if config.PanicLevel == 0 {
		config.PanicLevel = iface.Fatal
	}
}
