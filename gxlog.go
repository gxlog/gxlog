package gxlog

import (
	"fmt"
	"time"
)

const (
	cMapInitCap = 256
)

type locator struct {
	file string
	line int
}

type attribute struct {
	prefix       string
	contexts     []Context
	marked       bool
	countLimiter Filter
	timeLimiter  Filter
}

type Logger struct {
	*logger

	attr     attribute
	countMap map[locator]int64
	timeMap  map[locator]*timeQueue
}

func New(config *Config) *Logger {
	if config == nil {
		panic("nil config")
	}
	return &Logger{
		logger:   &logger{config: *config},
		countMap: make(map[locator]int64, cMapInitCap),
		timeMap:  make(map[locator]*timeQueue, cMapInitCap),
	}
}

func (this *Logger) WithPrefix(prefix string) *Logger {
	clone := *this
	clone.attr.prefix = prefix
	return &clone
}

func (this *Logger) WithContext(kvs ...interface{}) *Logger {
	clone := *this
	clone.attr.contexts = copyAppendContexts(clone.attr.contexts, kvs)
	return &clone
}

func (this *Logger) WithMark(ok bool) *Logger {
	clone := *this
	clone.attr.marked = ok
	return &clone
}

func (this *Logger) WithCountLimit(batch, limit int64) *Logger {
	clone := *this
	clone.attr.countLimiter = func(record *Record) bool {
		loc := locator{
			file: record.File,
			line: record.Line,
		}
		n := this.countMap[loc]
		this.countMap[loc]++
		if n%batch < limit {
			return true
		}
		return false
	}
	return &clone
}

func (this *Logger) WithTimeLimit(duration time.Duration, limit int) *Logger {
	clone := *this
	clone.attr.timeLimiter = func(record *Record) bool {
		loc := locator{
			file: record.File,
			line: record.Line,
		}
		queue := this.timeMap[loc]
		if queue == nil {
			queue = newTimeQueue(duration, limit)
			this.timeMap[loc] = queue
		}
		return queue.Enqueue(record.Time)
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
