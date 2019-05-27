package text

import (
	"github.com/gxlog/gxlog/iface"
)

type elementFormatter interface {
	FormatElement(buf []byte, record *iface.Record) []byte
}

var newFormatterFuncMap = map[string]func(property, fmtspec string) elementFormatter{
	"time":    newTimeFormatter,
	"level":   newLevelFormatter,
	"file":    newFileFormatter,
	"line":    newLineFormatter,
	"pkg":     newPkgFormatter,
	"func":    newFuncFormatter,
	"msg":     newMsgFormatter,
	"prefix":  newPrefixFormatter,
	"context": newContextFormatter,
}

type headerAppender struct {
	formatter  elementFormatter
	staticText string
}

func newHeaderAppender(element, property, fmtspec, staticText string) *headerAppender {
	newFunc := newFormatterFuncMap[element]
	if newFunc == nil {
		return nil
	}
	return &headerAppender{
		formatter:  newFunc(property, fmtspec),
		staticText: staticText,
	}
}

func (appender *headerAppender) AppendHeader(buf []byte, record *iface.Record) []byte {
	buf = append(buf, appender.staticText...)
	return appender.formatter.FormatElement(buf, record)
}
