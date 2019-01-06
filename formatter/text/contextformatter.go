package text

import (
	"fmt"
	"strings"

	"github.com/gxlog/gxlog/iface"
)

type contextFormatter struct {
	formatter func([]byte, []iface.Context) []byte
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

func (formatter *contextFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if formatter.fmtspec == "%s" {
		return formatter.formatter(buf, record.Aux.Contexts)
	}
	formatter.buf = formatter.buf[:0]
	formatter.buf = formatter.formatter(formatter.buf, record.Aux.Contexts)
	return append(buf, fmt.Sprintf(formatter.fmtspec, formatter.buf)...)
}

func selectFormatter(property string) func([]byte, []iface.Context) []byte {
	if strings.ToLower(property) == "list" {
		return formatList
	}
	return formatPair
}

func formatPair(buf []byte, contexts []iface.Context) []byte {
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

func formatList(buf []byte, contexts []iface.Context) []byte {
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
