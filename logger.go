package gxlog

import (
	"container/list"
	"fmt"
)

type link struct {
	ft Formatter
	wt Writer
}

type Logger struct {
	level    LogLevel
	links    list.List
	gatherer gatherer
}

func (this *Logger) Link(ft Formatter, wt Writer) bool {
	lnk := link{ft, wt}
	if this.linkExists(lnk) {
		return false
	}
	this.links.PushBack(lnk)
	return true
}

func (this *Logger) LinkBefore(ft Formatter, wt Writer, mark link) bool {
	lnk := link{ft, wt}
	if this.linkExists(lnk) {
		return false
	}
	for e := this.links.Front(); e != nil; e = e.Next() {
		if e.Value.(link) == mark {
			this.links.InsertBefore(lnk, e)
			return true
		}
	}
	return false
}

func (this *Logger) Unlink(ft Formatter, wt Writer) bool {
	lnk := link{ft, wt}
	for e := this.links.Front(); e != nil; e = e.Next() {
		if e.Value.(link) == lnk {
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
		lnk := e.Value.(link)
		formatter := lnk.ft
		format, ok := formatMap[formatter]
		if !ok {
			format = formatter.Format(record)
			formatMap[formatter] = format
		}
		lnk.wt.Write(format, record)
	}
}

func (this *Logger) linkExists(lnk link) bool {
	for e := this.links.Front(); e != nil; e = e.Next() {
		if e.Value.(link) == lnk {
			return true
		}
	}
	return false
}
