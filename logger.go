package gxlog

import (
	"fmt"
	"os"
)

const (
	cCallDepth = 4
)

type logger struct {
	linkSlots    [MaxLinkSlot]*link
	compactSlots []*link
	gatherer     gatherer
	level        LogLevel
	exitOnFatal  bool
}

func (this *logger) Log(calldepth int, level LogLevel, actions []Action, args []interface{}) {
	if this.level <= level {
		this.write(calldepth, level, actions, fmt.Sprint(args...))
	}
	if this.exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}
}

func (this *logger) Logf(calldepth int, level LogLevel, actions []Action,
	fmtstr string, args []interface{}) {
	if this.level <= level {
		this.write(calldepth, level, actions, fmt.Sprintf(fmtstr, args...))
	}
	if this.exitOnFatal && level == LevelFatal {
		os.Exit(1)
	}
}

func (this *logger) Panic(actions []Action, args []interface{}) {
	msg := fmt.Sprint(args...)
	if this.level <= LevelFatal {
		this.write(0, LevelFatal, actions, msg)
	}
	panic(msg)
}

func (this *logger) Panicf(actions []Action, fmtstr string, args []interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	if this.level <= LevelFatal {
		this.write(0, LevelFatal, actions, msg)
	}
	panic(msg)
}

func (this *logger) GetLevel() LogLevel {
	return this.level
}

func (this *logger) SetLevel(level LogLevel) {
	this.level = level
}

func (this *logger) GetExitOnFatal() bool {
	return this.exitOnFatal
}

func (this *logger) SetExitOnFatal(ok bool) {
	this.exitOnFatal = ok
}

func (this *logger) write(calldepth int, level LogLevel, actions []Action, msg string) {
	record := this.gatherer.Gather(calldepth+cCallDepth, level, msg)
	for _, action := range actions {
		action(record)
	}
	for _, lnk := range this.compactSlots {
		lnk.writer.Write(lnk.formatter.Format(record), record)
	}
}
