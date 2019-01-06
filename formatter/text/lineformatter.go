package text

import (
	"fmt"
	"strconv"

	"github.com/gxlog/gxlog/iface"
)

type lineFormatter struct {
	property string
	fmtspec  string
}

func newLineFormatter(property, fmtspec string) *lineFormatter {
	if fmtspec == "" {
		fmtspec = "%d"
	}
	return &lineFormatter{property: property, fmtspec: fmtspec}
}

func (formatter *lineFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if formatter.fmtspec == "%d" {
		return append(buf, strconv.Itoa(record.Line)...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, record.Line)...)
}
