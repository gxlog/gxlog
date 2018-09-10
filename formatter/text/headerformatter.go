package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05.000000000"
)

type headerFormatter interface {
	formatHeader(buf []byte, record *gxlog.Record) []byte
}

type staticFormatter struct {
	content []byte
}

func (this *staticFormatter) formatHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, this.content...)
}

type timeFormatter struct {
	property string
	fmtspec  string
}

func createTimeFormatter(property, fmtspec string) *timeFormatter {
	if property == "" {
		property = DefaultTimeLayout
	}
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &timeFormatter{property: property, fmtspec: fmtspec}
}

func (this *timeFormatter) formatHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Time.Format(this.property))...)
}

type levelFormatter struct {
	property string
	fmtspec  string
}

func createLevelFormatter(property, fmtspec string) *levelFormatter {
	if fmtspec == "" {
		fmtspec = "%-5s"
	}
	return &levelFormatter{property: property, fmtspec: fmtspec}
}

func (this *levelFormatter) formatHeader(buf []byte, record *gxlog.Record) []byte {
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

type pathnameFormatter struct {
	property string
	fmtspec  string
}

func createPathnameFormatter(property, fmtspec string) *pathnameFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &pathnameFormatter{property: property, fmtspec: fmtspec}
}

func (this *pathnameFormatter) formatHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Pathname)...)
}

type lineFormatter struct {
	property string
	fmtspec  string
}

func createLineFormatter(property, fmtspec string) *lineFormatter {
	if fmtspec == "" {
		fmtspec = "%d"
	}
	return &lineFormatter{property: property, fmtspec: fmtspec}
}

func (this *lineFormatter) formatHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Line)...)
}

type funcFormatter struct {
	property string
	fmtspec  string
}

func createFuncFormatter(property, fmtspec string) *funcFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &funcFormatter{property: property, fmtspec: fmtspec}
}

func (this *funcFormatter) formatHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Func)...)
}

type msgFormatter struct {
	property string
	fmtspec  string
}

func createMsgFormatter(property, fmtspec string) *msgFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &msgFormatter{property: property, fmtspec: fmtspec}
}

func (this *msgFormatter) formatHeader(buf []byte, record *gxlog.Record) []byte {
	return append(buf, fmt.Sprintf(this.fmtspec, record.Msg)...)
}
