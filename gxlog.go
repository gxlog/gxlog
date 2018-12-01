package gxlog

type Logger struct {
	*logger
}

func New() *Logger {
	return &Logger{&logger{}}
}
