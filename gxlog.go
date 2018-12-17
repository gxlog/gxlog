package gxlog

import (
	"fmt"
	"time"
)

const (
	cMapInitCap = 256
)

type Dynamic func(key interface{}) interface{}

type dynamicContext struct {
	key   interface{}
	value Dynamic
}

type locator struct {
	file string
	line int
}

type attribute struct {
	prefix          string
	contexts        []Context
	dynamicContexts []dynamicContext
	marked          bool
	countLimiter    Filter
	timeLimiter     Filter
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
	clone.appendContexts(kvs)
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

func (this *Logger) appendContexts(kvs []interface{}) {
	dynamicContexts := this.attr.dynamicContexts
	contexts := this.attr.contexts
	for len(kvs) >= 2 {
		dynamic, ok := kvs[1].(Dynamic)
		if ok {
			dynamicContexts = append(dynamicContexts, dynamicContext{
				key:   kvs[0],
				value: dynamic,
			})
		} else {
			contexts = append(contexts, Context{
				Key:   fmt.Sprint(kvs[0]),
				Value: fmt.Sprint(kvs[1]),
			})
		}
		kvs = kvs[2:]
	}
	// slicing to set capacity to length, force next appending to reallocate memory
	this.attr.contexts = contexts[:len(contexts):len(contexts)]
	this.attr.dynamicContexts = dynamicContexts[:len(dynamicContexts):len(dynamicContexts)]
}
