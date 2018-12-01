package gxlog

import (
	"fmt"
)

type Action func(*Record)

type Logger struct {
	*logger
	actions []Action
}

func New(config *Config) *Logger {
	return &Logger{
		logger: &logger{
			level:       config.Level,
			exitOnFatal: config.ExitOnFatal,
		},
	}
}

func (this *Logger) WithPrefix(prefix string) *Logger {
	actions := copyAppend(this.actions, func(record *Record) {
		record.Prefix = prefix
	})
	return &Logger{
		logger:  this.logger,
		actions: actions,
	}
}

func (this *Logger) WithContext(kvs ...interface{}) *Logger {
	contexts := make([]Context, 0, len(kvs)/2)
	for len(kvs) >= 2 {
		key := fmt.Sprint(kvs[0])
		value := fmt.Sprint(kvs[1])
		contexts = append(contexts, Context{key, value})
		kvs = kvs[2:]
	}
	actions := copyAppend(this.actions, func(record *Record) {
		record.Contexts = append(record.Contexts, contexts...)
	})
	return &Logger{
		logger:  this.logger,
		actions: actions,
	}
}

func (this *Logger) WithMark() *Logger {
	actions := copyAppend(this.actions, func(record *Record) {
		record.Marked = true
	})
	return &Logger{
		logger:  this.logger,
		actions: actions,
	}
}

func copyAppend(actions []Action, action Action) []Action {
	newActions := make([]Action, 0, len(actions)+1)
	newActions = append(newActions, actions...)
	newActions = append(newActions, action)
	return newActions
}
