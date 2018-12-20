package gxlog

import (
	"errors"
	"fmt"
)

func (this *Logger) Trace(args ...interface{}) {
	this.logger.Log(0, LevelTrace, &this.attr, args...)
}

func (this *Logger) Tracef(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelTrace, &this.attr, fmtstr, args...)
}

func (this *Logger) Debug(args ...interface{}) {
	this.logger.Log(0, LevelDebug, &this.attr, args...)
}

func (this *Logger) Debugf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelDebug, &this.attr, fmtstr, args...)
}

func (this *Logger) Info(args ...interface{}) {
	this.logger.Log(0, LevelInfo, &this.attr, args...)
}

func (this *Logger) Infof(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelInfo, &this.attr, fmtstr, args...)
}

func (this *Logger) Warn(args ...interface{}) {
	this.logger.Log(0, LevelWarn, &this.attr, args...)
}

func (this *Logger) Warnf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelWarn, &this.attr, fmtstr, args...)
}

func (this *Logger) Error(args ...interface{}) {
	this.logger.Log(0, LevelError, &this.attr, args...)
}

func (this *Logger) Errorf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelError, &this.attr, fmtstr, args...)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.logger.Log(0, LevelFatal, &this.attr, args...)
}

func (this *Logger) Fatalf(fmtstr string, args ...interface{}) {
	this.logger.Logf(0, LevelFatal, &this.attr, fmtstr, args...)
}

func (this *Logger) Panic(args ...interface{}) {
	this.logger.Panic(&this.attr, args...)
}

func (this *Logger) Panicf(fmtstr string, args ...interface{}) {
	this.logger.Panicf(&this.attr, fmtstr, args...)
}

func (this *Logger) Time(args ...interface{}) func() {
	return this.logger.Time(&this.attr, args...)
}

func (this *Logger) Timef(fmtstr string, args ...interface{}) func() {
	return this.logger.Timef(&this.attr, fmtstr, args...)
}

func (this *Logger) Log(calldepth int, level Level, args ...interface{}) {
	this.logger.Log(calldepth, level, &this.attr, args...)
}

func (this *Logger) Logf(calldepth int, level Level, fmtstr string, args ...interface{}) {
	this.logger.Logf(calldepth, level, &this.attr, fmtstr, args...)
}

func (this *Logger) LogError(level Level, text string) error {
	this.logger.Log(0, level, &this.attr, text)
	return errors.New(text)
}

func (this *Logger) LogErrorf(level Level, fmtstr string, args ...interface{}) error {
	err := fmt.Errorf(fmtstr, args...)
	this.logger.Log(0, level, &this.attr, err.Error())
	return err
}
