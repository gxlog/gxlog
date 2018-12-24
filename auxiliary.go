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
	Key   interface{}
	Value Dynamic
}

type locator struct {
	File string
	Line int
}

type copyOnWrite struct {
	Prefix          string
	Contexts        []Context
	DynamicContexts []dynamicContext
	Marked          bool
	CountLimiter    Filter
	TimeLimiter     Filter
}

func (this *Logger) WithPrefix(prefix string) *Logger {
	clone := *this
	clone.attr.Prefix = prefix
	return &clone
}

func (this *Logger) WithContext(kvs ...interface{}) *Logger {
	clone := *this
	clone.attr.Contexts, clone.attr.DynamicContexts =
		appendContexts(clone.attr.Contexts, clone.attr.DynamicContexts, kvs)
	return &clone
}

func (this *Logger) WithMark(ok bool) *Logger {
	clone := *this
	clone.attr.Marked = ok
	return &clone
}

func (this *Logger) WithCountLimit(batch, limit int64) *Logger {
	clone := *this
	clone.attr.CountLimiter = func(record *Record) bool {
		loc := locator{
			File: record.File,
			Line: record.Line,
		}
		n := this.countMap[loc]
		this.countMap[loc]++
		return n%batch < limit
	}
	return &clone
}

func (this *Logger) WithTimeLimit(duration time.Duration, limit int) *Logger {
	clone := *this
	clone.attr.TimeLimiter = func(record *Record) bool {
		loc := locator{
			File: record.File,
			Line: record.Line,
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

func appendContexts(contexts []Context, dynamicContexts []dynamicContext,
	kvs []interface{}) ([]Context, []dynamicContext) {
	for len(kvs) >= 2 {
		dynamic, ok := kvs[1].(Dynamic)
		if ok {
			dynamicContexts = append(dynamicContexts, dynamicContext{
				Key:   kvs[0],
				Value: dynamic,
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
	return contexts[:len(contexts):len(contexts)],
		dynamicContexts[:len(dynamicContexts):len(dynamicContexts)]
}
