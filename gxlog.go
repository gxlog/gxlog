package gxlog

type Logger struct {
	*logger
	actions []func(*Record)
}

func New() *Logger {
	return &Logger{
		logger: &logger{},
	}
}

func (this *Logger) WithPrefix(prefix string) *Logger {
	actions := make([]func(*Record), 0, len(this.actions)+1)
	actions = append(actions, this.actions...)
	actions = append(actions, func(record *Record) {
		record.Prefix = prefix
	})
	return &Logger{
		logger:  this.logger,
		actions: actions,
	}
}
