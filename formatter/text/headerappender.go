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

type staticAppender struct {
	content []byte
}

func (this *staticAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, this.content...)
}

type timeAppender struct {
	property string
	fmtspec  string
}

func createTimeAppender(property, fmtspec string) *timeAppender {
	if property == "" {
		property = DefaultTimeLayout
	}
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &timeAppender{property: property, fmtspec: fmtspec}
}

func (this *timeAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Time.Format(this.property))...)
}

type levelAppender struct {
	property string
	fmtspec  string
}

func createLevelAppender(property, fmtspec string) *levelAppender {
	if fmtspec == "" {
		fmtspec = "%-5s"
	}
	return &levelAppender{property: property, fmtspec: fmtspec}
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
	return append(buf, fmt.Sprintf(this.fmtspec, level)...)
}

type pathnameAppender struct {
	property string
	fmtspec  string
}

func createPathnameAppender(property, fmtspec string) *pathnameAppender {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &pathnameAppender{property: property, fmtspec: fmtspec}
}

func (this *pathnameAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Pathname)...)
}

type lineAppender struct {
	property string
	fmtspec  string
}

func createLineAppender(property, fmtspec string) *lineAppender {
	if fmtspec == "" {
		fmtspec = "%d"
	}
	return &lineAppender{property: property, fmtspec: fmtspec}
}

func (this *lineAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Line)...)
}

type funcAppender struct {
	property string
	fmtspec  string
}

func createFuncAppender(property, fmtspec string) *funcAppender {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &funcAppender{property: property, fmtspec: fmtspec}
}

func (this *funcAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Func)...)
}

type msgAppender struct {
	property string
	fmtspec  string
}

func createMsgAppender(property, fmtspec string) *msgAppender {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &msgAppender{property: property, fmtspec: fmtspec}
}

func (this *msgAppender) appendHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Msg)...)
}
