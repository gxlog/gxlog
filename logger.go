package gxlog

import (
	"container/list"
)

type link struct {
	f Formatter
	w Writer
}

type Logger struct {
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

func (this *Logger) UnlinkAll(f Formatter, w Writer) {
	this.links.Init()
}

func (this *Logger) Log() {
	formatMap := make(map[Formatter][]byte)
	i := this.gatherer.gather()
	for e := this.links.Front(); e != nil; e = e.Next() {
		l := e.Value.(link)
		formatter := l.f
		format, ok := formatMap[formatter]
		if !ok {
			format = formatter.Format(i)
			formatMap[formatter] = format
		}
		l.w.Write(format)
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
