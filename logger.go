package gxlog

import (
	"fmt"
)

type Link struct {
	FT Formatter
	WT Writer
}

type Logger struct {
	level    LogLevel
	links    []Link
	gatherer gatherer
}

func (this *Logger) Link(ft Formatter, wt Writer) {
	this.links = append(this.links, Link{ft, wt})
}

func (this *Logger) LinkAll(links []Link) {
	this.links = append(this.links, links...)
}

func (this *Logger) UnlinkAll() {
	this.links = nil
}

func (this *Logger) ResetAll(links []Link) {
	this.links = make([]Link, len(links))
	copy(this.links, links)
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
	if this.level <= level {
		this.write(level, fmt.Sprint(args...))
	}
}

func (this *Logger) Logf(level LogLevel, fmtstr string, args []interface{}) {
	if this.level <= level {
		this.write(level, fmt.Sprintf(fmtstr, args...))
	}
}

func (this *Logger) write(level LogLevel, msg string) {
	record := this.gatherer.gather(4, level, msg)
	for _, ln := range this.links {
		ln.WT.Write(ln.FT.Format(record), record)
	}
}
