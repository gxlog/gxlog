package gxlog

import (
	"container/list"
	"fmt"
)

type link struct {
	f Formatter
	w Writer
}

type Logger struct {
	level    LogLevel
	links    list.List
	gatherer gatherer
}

func (this *Logger) Link(f Formatter, w Writer) bool {
	l := link{f, w}
	if this.linkExists(l) {
		return false
	}
	this.links.PushBack(l)
	return true
}

func (this *Logger) LinkBefore(f Formatter, w Writer, mark link) bool {
	l := link{f, w}
	if this.linkExists(l) {
		return false
	}
	for e := this.links.Front(); e != nil; e = e.Next() {
		if e.Value.(link) == mark {
			this.links.InsertBefore(l, e)
			return true
		}
	}
	return false
}

func (this *Logger) Unlink(f Formatter, w Writer) bool {
	l := link{f, w}
	for e := this.links.Front(); e != nil; e = e.Next() {
		if e.Value.(link) == l {
			this.links.Remove(e)
			return true
		}
	}
	return false
}

func (this *Logger) UnlinkAll() {
	this.links.Init()
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
	formatMap := make(map[Formatter][]byte)
	record := this.gatherer.gather(4, level, msg)
	for e := this.links.Front(); e != nil; e = e.Next() {
		l := e.Value.(link)
		formatter := l.f
		format, ok := formatMap[formatter]
		if !ok {
			format = formatter.Format(record)
			formatMap[formatter] = format
		}
		l.w.Write(format, record)
	}
}

func (this *Logger) linkExists(l link) bool {
	for e := this.links.Front(); e != nil; e = e.Next() {
		if e.Value.(link) == l {
			return true
		}
	}
	return false
}
