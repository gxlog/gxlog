package text

import (
	"fmt"
	"strings"

	"github.com/gxlog/gxlog"
)

type contextFormatter struct {
	formatter func([]byte, []gxlog.Context) []byte
	fmtspec   string
	buf       []byte
}

func newContextFormatter(property, fmtspec string) *contextFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &contextFormatter{
		formatter: selectFormatter(property),
		fmtspec:   fmtspec,
	}
}

func (this *contextFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if this.fmtspec == "%s" {
		return this.formatter(buf, record.Aux.Contexts)
	} else {
		this.buf = this.buf[:0]
		this.buf = this.formatter(this.buf, record.Aux.Contexts)
		return append(buf, fmt.Sprintf(this.fmtspec, this.buf)...)
	}
}

func selectFormatter(property string) func([]byte, []gxlog.Context) []byte {
	if strings.ToLower(property) == "list" {
		return formatList
	}
	return formatPair
}

func formatPair(buf []byte, contexts []gxlog.Context) []byte {
	left := "("
	for _, ctx := range contexts {
		buf = append(buf, left...)
		buf = append(buf, ctx.Key...)
		buf = append(buf, ": "...)
		buf = append(buf, ctx.Value...)
		buf = append(buf, ')')
		left = " ("
	}
	return buf
}

func formatList(buf []byte, contexts []gxlog.Context) []byte {
	begin := ""
	for _, ctx := range contexts {
		buf = append(buf, begin...)
		buf = append(buf, ctx.Key...)
		buf = append(buf, ": "...)
		buf = append(buf, ctx.Value...)
		begin = ", "
	}
	return buf
}
