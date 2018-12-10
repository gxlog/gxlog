package gxlog

func (this *Logger) Trace(args ...interface{}) {
	this.logger.Log(0, LevelTrace, &this.aux, args)
}

func (this *Logger) Tracef(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelTrace, &this.aux, fmtstr, args)
}

func (this *Logger) Debug(args ...interface{}) {
	this.logger.Log(0, LevelDebug, &this.aux, args)
}

func (this *Logger) Debugf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelDebug, &this.aux, fmtstr, args)
}

func (this *Logger) Info(args ...interface{}) {
	this.logger.Log(0, LevelInfo, &this.aux, args)
}

func (this *Logger) Infof(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelInfo, &this.aux, fmtstr, args)
}

func (this *Logger) Warn(args ...interface{}) {
	this.logger.Log(0, LevelWarn, &this.aux, args)
}

func (this *Logger) Warnf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelWarn, &this.aux, fmtstr, args)
}

func (this *Logger) Error(args ...interface{}) {
	this.logger.Log(0, LevelError, &this.aux, args)
}

func (this *Logger) Errorf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelError, &this.aux, fmtstr, args)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.logger.Log(0, LevelFatal, &this.aux, args)
}

func (this *Logger) Fatalf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelFatal, &this.aux, fmtstr, args)
}

func (this *Logger) Panic(args ...interface{}) {
	this.logger.Panic(&this.aux, args)
}

func (this *Logger) Panicf(fmtstr string, args ...interface{}) {
	this.logger.Panicf(&this.aux, fmtstr, args)
}

func (this *Logger) Time(args ...interface{}) func() {
	return this.logger.Time(&this.aux, args)
}

func (this *Logger) Timef(fmtstr string, args ...interface{}) func() {
	return this.logger.Timef(&this.aux, fmtstr, args)
}

func (this *Logger) Log(calldepth int, level Level, args ...interface{}) {
	this.logger.Log(calldepth, level, &this.aux, args)
}

func (this *Logger) Logf(calldepth int, level Level, fmtstr string, args ...interface{}) {
	this.logger.Logf(calldepth, level, &this.aux, fmtstr, args)
}
