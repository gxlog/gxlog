package text

import (
	"fmt"

	"github.com/gxlog/gxlog/iface"
)

type prefixFormatter struct {
	property string
	fmtspec  string
}

func newPrefixFormatter(property, fmtspec string) *prefixFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &prefixFormatter{property: property, fmtspec: fmtspec}
}

func (formatter *prefixFormatter) FormatElement(buf []byte,
	record *iface.Record) []byte {
	if formatter.fmtspec == "%s" {
		return append(buf, record.Aux.Prefix...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, record.Aux.Prefix)...)
}
