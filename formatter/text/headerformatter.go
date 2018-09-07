package text

import (
	"bytes"
	"fmt"

	"github.com/gratonos/gxlog"
)

const (
	DefaultTimeLayout = "2006-01-02 15:04:05.000000000"
)

type headerFormatter interface {
	formatHeader(*gxlog.Record) []byte
}

type staticFormatter struct {
	content []byte
}

func (this *staticFormatter) formatHeader(record *gxlog.Record) []byte {
	return this.content
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

func (this *timeFormatter) formatHeader(record *gxlog.Record) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, this.fmtspec, record.Time.Format(this.property))
	return buf.Bytes()
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

func (this *levelFormatter) formatHeader(record *gxlog.Record) []byte {
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
	var buf bytes.Buffer
	fmt.Fprintf(&buf, this.fmtspec, level)
	return buf.Bytes()
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

func (this *pathnameFormatter) formatHeader(record *gxlog.Record) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, this.fmtspec, record.Pathname)
	return buf.Bytes()
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

func (this *lineFormatter) formatHeader(record *gxlog.Record) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, this.fmtspec, record.Line)
	return buf.Bytes()
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

func (this *funcFormatter) formatHeader(record *gxlog.Record) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, this.fmtspec, record.Func)
	return buf.Bytes()
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

func (this *msgFormatter) formatHeader(record *gxlog.Record) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, this.fmtspec, record.Msg)
	return buf.Bytes()
}
