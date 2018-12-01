package gxlog

func (this *Logger) Trace(args ...interface{}) {
	this.Log(LevelTrace, args)
}

func (this *Logger) Tracef(fmtstr string, args ...interface{}) {
	this.Logf(LevelTrace, fmtstr, args)
}

func (this *Logger) Debug(args ...interface{}) {
	this.Log(LevelDebug, args)
}

func (this *Logger) Debugf(fmtstr string, args ...interface{}) {
	this.Logf(LevelDebug, fmtstr, args)
}

func (this *Logger) Info(args ...interface{}) {
	this.Log(LevelInfo, args)
}

func (this *Logger) Infof(fmtstr string, args ...interface{}) {
	this.Logf(LevelInfo, fmtstr, args)
}

func (this *Logger) Warn(args ...interface{}) {
	this.Log(LevelWarn, args)
}

func (this *Logger) Warnf(fmtstr string, args ...interface{}) {
	this.Logf(LevelWarn, fmtstr, args)
}

func (this *Logger) Error(args ...interface{}) {
	this.Log(LevelError, args)
}

func (this *Logger) Errorf(fmtstr string, args ...interface{}) {
	this.Logf(LevelError, fmtstr, args)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.Log(LevelFatal, args)
}

func (this *Logger) Fatalf(fmtstr string, args ...interface{}) {
	this.Logf(LevelFatal, fmtstr, args)
}

func (this *Logger) Log(level LogLevel, args []interface{}) {
	this.logger.Log(level, this.actions, args)
}

func (this *Logger) Logf(level LogLevel, fmtstr string, args []interface{}) {
	this.logger.Logf(level, this.actions, fmtstr, args)
}
