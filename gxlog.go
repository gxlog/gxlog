package gxlog

import (
	"fmt"
)

const (
	cMapInitCap = 256
)

type locator struct {
	file string
	line int
}

type attribute struct {
	aux          Auxiliary
	countLimiter Filter
}

type Logger struct {
	*logger

	attr attribute
}

func New(config *Config) *Logger {
	if config == nil {
		panic("nil config")
	}
	return &Logger{
		logger: &logger{
			level:      config.Level,
			trackLevel: config.TrackLevel,
			exitLevel:  config.ExitLevel,
			filter:     config.Filter,
			limit:      config.Limit,
			countMap:   make(map[locator]int64, cMapInitCap),
		},
	}
}

func (this *Logger) WithPrefix(prefix string) *Logger {
	clone := *this
	clone.attr.aux.Prefix = prefix
	return &clone
}

func (this *Logger) WithContext(kvs ...interface{}) *Logger {
	clone := *this
	clone.attr.aux.Contexts = copyAppendContexts(clone.attr.aux.Contexts, kvs)
	return &clone
}

func (this *Logger) WithMark(ok bool) *Logger {
	clone := *this
	clone.attr.aux.Marked = ok
	return &clone
}

func (this *Logger) WithCountLimit(batch, limit int64) *Logger {
	clone := *this
	clone.attr.countLimiter = func(record *Record) bool {
		loc := locator{
			file: record.File,
			line: record.Line,
		}
		n := this.logger.countMap[loc]
		this.logger.countMap[loc]++
		if n%batch < limit {
			return true
		}
		return false
	}
	return &clone
}

func copyAppendContexts(dst []Context, kvs []interface{}) []Context {
	contexts := make([]Context, 0, len(dst)+len(kvs)/2)
	contexts = append(contexts, dst...)
	for len(kvs) >= 2 {
		key := fmt.Sprint(kvs[0])
		value := fmt.Sprint(kvs[1])
		contexts = append(contexts, Context{Key: key, Value: value})
		kvs = kvs[2:]
	}
	return contexts
}
