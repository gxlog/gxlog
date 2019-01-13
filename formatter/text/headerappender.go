package text

import (
	"github.com/gxlog/gxlog/iface"
)

type elementFormatter interface {
	FormatElement(buf []byte, record *iface.Record) []byte
}

type headerAppender struct {
	formatter  elementFormatter
	staticText string
}

func newHeaderAppender(element, property, fmtspec,
	staticText string) *headerAppender {
	var formatter elementFormatter
	switch element {
	case "time":
		formatter = newTimeFormatter(property, fmtspec)
	case "level":
		formatter = newLevelFormatter(property, fmtspec)
	case "file":
		formatter = newFileFormatter(property, fmtspec)
	case "line":
		formatter = newLineFormatter(property, fmtspec)
	case "pkg":
		formatter = newPkgFormatter(property, fmtspec)
	case "func":
		formatter = newFuncFormatter(property, fmtspec)
	case "msg":
		formatter = newMsgFormatter(property, fmtspec)
	case "prefix":
		formatter = newPrefixFormatter(property, fmtspec)
	case "context":
		formatter = newContextFormatter(property, fmtspec)
	}
	if formatter != nil {
		return &headerAppender{formatter: formatter, staticText: staticText}
	}
	return nil
}

func (appender *headerAppender) AppendHeader(buf []byte,
	record *iface.Record) []byte {
	buf = append(buf, appender.staticText...)
	return appender.formatter.FormatElement(buf, record)
}
