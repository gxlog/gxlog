package gxlog

func (this *Logger) Trace(args ...interface{}) {
	this.logger.Log(0, LevelTrace, this.actions, args)
}

func (this *Logger) Tracef(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelTrace, this.actions, fmtstr, args)
}

func (this *Logger) Debug(args ...interface{}) {
	this.logger.Log(0, LevelDebug, this.actions, args)
}

func (this *Logger) Debugf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelDebug, this.actions, fmtstr, args)
}

func (this *Logger) Info(args ...interface{}) {
	this.logger.Log(0, LevelInfo, this.actions, args)
}

func (this *Logger) Infof(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelInfo, this.actions, fmtstr, args)
}

func (this *Logger) Warn(args ...interface{}) {
	this.logger.Log(0, LevelWarn, this.actions, args)
}

func (this *Logger) Warnf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelWarn, this.actions, fmtstr, args)
}

func (this *Logger) Error(args ...interface{}) {
	this.logger.Log(0, LevelError, this.actions, args)
}

func (this *Logger) Errorf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelError, this.actions, fmtstr, args)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.logger.Log(0, LevelFatal, this.actions, args)
}

func (this *Logger) Fatalf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelFatal, this.actions, fmtstr, args)
}

func (this *Logger) Panic(args ...interface{}) {
	this.logger.Panic(this.actions, args)
}

func (this *Logger) Panicf(fmtstr string, args ...interface{}) {
	this.logger.Panicf(this.actions, fmtstr, args)
}

func (this *Logger) Log(calldepth int, level LogLevel, args ...interface{}) {
	this.logger.Log(calldepth, level, this.actions, args)
}

func (this *Logger) Logf(calldepth int, level LogLevel, fmtstr string, args ...interface{}) {
	this.logger.Logf(calldepth, level, this.actions, fmtstr, args)
}
