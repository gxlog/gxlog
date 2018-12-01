package gxlog

import (
	"fmt"
)

const (
	cCallDepth = 5
)

type logger struct {
	linkSlots    [MaxLinkSlot]*link
	compactSlots []*link
	level        LogLevel
	gatherer     gatherer
}

func (this *logger) Log(level LogLevel, args []interface{}) {
	if this.level <= level {
		this.write(level, fmt.Sprint(args...))
	}
}

func (this *logger) Logf(level LogLevel, fmtstr string, args []interface{}) {
	if this.level <= level {
		this.write(level, fmt.Sprintf(fmtstr, args...))
	}
}

func (this *logger) write(level LogLevel, msg string) {
	record := this.gatherer.Gather(cCallDepth, level, msg)
	for _, lnk := range this.compactSlots {
		lnk.writer.Write(lnk.formatter.Format(record), record)
	}
}
