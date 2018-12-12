package text

import (
	"fmt"

	"github.com/gxlog/gxlog"
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

func (this *prefixFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if this.fmtspec == "%s" {
		return append(buf, record.Aux.Prefix...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, record.Aux.Prefix)...)
	}
}
