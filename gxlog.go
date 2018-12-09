package gxlog

import (
	"fmt"
)

type Logger struct {
	*logger

	aux Auxiliary
}

func New(config *Config) *Logger {
	if config == nil {
		panic("nil config")
	}
	return &Logger{
		logger: &logger{
			level:       config.Level,
			filter:      config.Filter,
			exitOnFatal: config.ExitOnFatal,
		},
	}
}

func (this *Logger) WithPrefix(prefix string) *Logger {
	clone := *this
	clone.aux.Prefix = prefix
	return &clone
}

func (this *Logger) WithContext(kvs ...interface{}) *Logger {
	clone := *this
	clone.aux.Contexts = copyAppendContexts(clone.aux.Contexts, kvs)
	return &clone
}

func (this *Logger) WithMark(ok bool) *Logger {
	clone := *this
	clone.aux.Marked = ok
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
