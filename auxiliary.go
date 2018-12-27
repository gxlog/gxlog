package gxlog

import (
	"fmt"
	"time"
)

const mapInitCap = 256

// The Dynamic type defines a function type. A value of Dynamic will be regarded
// as the value getter of a dynamic context key-value pair when it is as an
// argument to (*Logger).WithContext.
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

// WithPrefix returns a new Logger that is a shallow copy of the caller.
// With the new Logger, all the logs it outputs will have the prefix as
// long as the Prefix flag of the Logger is set.
func (log *Logger) WithPrefix(prefix string) *Logger {
	clone := *log
	clone.attr.Prefix = prefix
	return &clone
}

// WithContext returns a new Logger that is a shallow copy of the caller.
// With the new Logger, all the logs it outputs will have all contexts
// concatenated as long as the Contexts flag of the Logger is set.
//
// The kvs is regarded as an interleaved key-value sequence,
// e.g. key1, value1, key2, value2 ...
// If the count of arguments is odd, the last argument will be ignored.
//
// WithContext also supports for dynamic context. If a value of type
// Dynamic is as the value of a key-value pair passed to WithContext,
// it will be regarded as the value getter of a dynamic context key-value
// pair. The value getter will be called whenever a new log emits.
// All the key-value pairs of dynamic contexts will be concatenated to the
// end of static contexts.
//
// ATTENTION: you should be very careful to concurrency safety or deadlocks
// with dynamic contexts.
func (log *Logger) WithContext(kvs ...interface{}) *Logger {
	clone := *log
	clone.attr.Contexts, clone.attr.DynamicContexts =
		appendContexts(clone.attr.Contexts, clone.attr.DynamicContexts, kvs)
	return &clone
}

// WithMark returns a new Logger that is a shallow copy of the caller.
// With the new Logger, all the logs it outputs will be marked as long
// as the Mark flag of the Logger is set.
func (log *Logger) WithMark(ok bool) *Logger {
	clone := *log
	clone.attr.Marked = ok
	return &clone
}

// WithCountLimit returns a new Logger that is a shallow copy of the caller.
// With the new Logger, the count of logs it outputs will be limited as long
// as the Limit flag of the Logger is set.
//
// The batch must be positive and the limit must NOT be negative. As a result,
// limit logs will be output every batch logs.
//
// THINK TWICE before you limit the output count of logs, you may miss logs
// which you need.
func (log *Logger) WithCountLimit(batch, limit int64) *Logger {
	if batch <= 0 {
		panic("gxlog.WithCountLimit: batch must be positive")
	}
	if limit < 0 {
		panic("gxlog.WithCountLimit: negative limit")
	}
	clone := *log
	clone.attr.CountLimiter = func(record *Record) bool {
		loc := locator{
			File: record.File,
			Line: record.Line,
		}
		n := log.countMap[loc]
		log.countMap[loc]++
		return n%batch < limit
	}
	return &clone
}

// WithTimeLimit returns a new Logger that is a shallow copy of the caller.
// With the new Logger, the count of logs it outputs will be limited as long
// as the Limit flag of the Logger is set.
//
// The duration must be positive and the limit must NOT be negative. As a result,
// at most limit logs will be output during any interval of duration.
//
// THINK TWICE before you limit the output count of logs, you may miss logs
// which you need.
//
// NOTICE: The space complexity of WithTimeLimit is O(limit). Try to specify
// reasonable duration and limit.
func (log *Logger) WithTimeLimit(duration time.Duration, limit int) *Logger {
	if duration <= 0 {
		panic("gxlog.WithTimeLimit: duration must be positive")
	}
	if limit < 0 {
		panic("gxlog.WithTimeLimit: negative limit")
	}
	clone := *log
	clone.attr.TimeLimiter = func(record *Record) bool {
		loc := locator{
			File: record.File,
			Line: record.Line,
		}
		queue := log.timeMap[loc]
		if queue == nil {
			queue = newTimeQueue(duration, limit)
			log.timeMap[loc] = queue
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
