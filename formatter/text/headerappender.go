package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05.000000"
)

type headerAppender interface {
	appendHeader(buf []byte, record *gxlog.Record) []byte
}

type timeAppender struct {
	property string
	fmtspec  string
	prefix   []byte
}

func createTimeAppender(property, fmtspec string, prefix []byte) *timeAppender {
	if property == "" {
		property = DefaultTimeLayout
	}
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &timeAppender{property: property, fmtspec: fmtspec, prefix: cloneBytes(prefix)}
}

func (this *timeAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	buf = append(buf, this.prefix...)
	return append(buf, fmt.Sprintf(this.fmtspec, record.Time.Format(this.property))...)
}

type levelAppender struct {
	property string
	fmtspec  string
	prefix   []byte
}

func createLevelAppender(property, fmtspec string, prefix []byte) *levelAppender {
	if fmtspec == "" {
		fmtspec = "%-5s"
	}
	return &levelAppender{property: property, fmtspec: fmtspec, prefix: cloneBytes(prefix)}
}

func (this *levelAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	var level string
	switch record.Level {
	case gxlog.LevelDebug:
		level = "DEBUG"
	case gxlog.LevelInfo:
		level = "INFO"
	case gxlog.LevelWarn:
		level = "WARN"
	case gxlog.LevelError:
		level = "ERROR"
	case gxlog.LevelFatal:
		level = "FATAL"
	}
	buf = append(buf, this.prefix...)
	return append(buf, fmt.Sprintf(this.fmtspec, level)...)
}

type pathnameAppender struct {
	property string
	fmtspec  string
	prefix   []byte
}

func createPathnameAppender(property, fmtspec string, prefix []byte) *pathnameAppender {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &pathnameAppender{property: property, fmtspec: fmtspec, prefix: cloneBytes(prefix)}
}

func (this *pathnameAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	buf = append(buf, this.prefix...)
	return append(buf, fmt.Sprintf(this.fmtspec, record.Pathname)...)
}

type lineAppender struct {
	property string
	fmtspec  string
	prefix   []byte
}

func createLineAppender(property, fmtspec string, prefix []byte) *lineAppender {
	if fmtspec == "" {
		fmtspec = "%d"
	}
	return &lineAppender{property: property, fmtspec: fmtspec, prefix: cloneBytes(prefix)}
}

func (this *lineAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	buf = append(buf, this.prefix...)
	return append(buf, fmt.Sprintf(this.fmtspec, record.Line)...)
}

type funcAppender struct {
	property string
	fmtspec  string
	prefix   []byte
}

func createFuncAppender(property, fmtspec string, prefix []byte) *funcAppender {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &funcAppender{property: property, fmtspec: fmtspec, prefix: cloneBytes(prefix)}
}

func (this *funcAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	buf = append(buf, this.prefix...)
	return append(buf, fmt.Sprintf(this.fmtspec, record.Func)...)
}

type msgAppender struct {
	property string
	fmtspec  string
	prefix   []byte
}

func createMsgAppender(property, fmtspec string, prefix []byte) *msgAppender {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &msgAppender{property: property, fmtspec: fmtspec, prefix: cloneBytes(prefix)}
}

func (this *msgAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	buf = append(buf, this.prefix...)
	return append(buf, fmt.Sprintf(this.fmtspec, record.Msg)...)
}
