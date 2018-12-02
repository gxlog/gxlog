package gxlog

func (this *Logger) Trace(args ...interface{}) {
	this.logger.Log(LevelTrace, this.actions, args)
}

func (this *Logger) Tracef(fmtstr string, args ...interface{}) {
	this.logger.Logf(LevelTrace, this.actions, fmtstr, args)
}

func (this *Logger) Debug(args ...interface{}) {
	this.logger.Log(LevelDebug, this.actions, args)
}

func (this *Logger) Debugf(fmtstr string, args ...interface{}) {
	this.logger.Logf(LevelDebug, this.actions, fmtstr, args)
}

func (this *Logger) Info(args ...interface{}) {
	this.logger.Log(LevelInfo, this.actions, args)
}

func (this *Logger) Infof(fmtstr string, args ...interface{}) {
	this.logger.Logf(LevelInfo, this.actions, fmtstr, args)
}

func (this *Logger) Warn(args ...interface{}) {
	this.logger.Log(LevelWarn, this.actions, args)
}

func (this *Logger) Warnf(fmtstr string, args ...interface{}) {
	this.logger.Logf(LevelWarn, this.actions, fmtstr, args)
}

func (this *Logger) Error(args ...interface{}) {
	this.logger.Log(LevelError, this.actions, args)
}

func (this *Logger) Errorf(fmtstr string, args ...interface{}) {
	this.logger.Logf(LevelError, this.actions, fmtstr, args)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.logger.Log(LevelFatal, this.actions, args)
}

func (this *Logger) Fatalf(fmtstr string, args ...interface{}) {
	this.logger.Logf(LevelFatal, this.actions, fmtstr, args)
}

func (this *Logger) Panic(args ...interface{}) {
	this.logger.Panic(this.actions, args)
}

func (this *Logger) Panicf(fmtstr string, args ...interface{}) {
	this.logger.Panicf(this.actions, fmtstr, args)
}

func (this *Logger) Log(level LogLevel, args ...interface{}) {
	this.logger.Log(level, this.actions, args)
}

func (this *Logger) Logf(level LogLevel, fmtstr string, args ...interface{}) {
	this.logger.Logf(level, this.actions, fmtstr, args)
}
