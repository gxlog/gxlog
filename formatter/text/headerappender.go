package text

import (
	"github.com/gratonos/gxlog"
)

type elementFormatter interface {
	FormatElement(record *gxlog.Record) string
}

type headerAppender struct {
	formatter elementFormatter
	prefix    []byte
}

func newHeaderAppender(element, property, fmtspec, prefix string) *headerAppender {
	var formatter elementFormatter
	switch element {
	case "time":
		formatter = newTimeFormatter(property, fmtspec)
	case "level":
		formatter = newLevelFormatter(property, fmtspec)
	case "pathname":
		formatter = newPathnameFormatter(property, fmtspec)
	case "line":
		formatter = newLineFormatter(property, fmtspec)
	case "func":
		formatter = newFuncFormatter(property, fmtspec)
	case "msg":
		formatter = newMsgFormatter(property, fmtspec)
	}
	if formatter != nil {
		return &headerAppender{formatter: formatter, prefix: []byte(prefix)}
	}
	return nil
}

func (this *headerAppender) AppendHeader(buf []byte, record *gxlog.Record) []byte {
	str := this.formatter.FormatElement(record)
	return append(append(buf, this.prefix...), str...)
}
