package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type pathnameFormatter struct {
	property string
	fmtspec  string
}

func newPathnameFormatter(property, fmtspec string) *pathnameFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &pathnameFormatter{property: property, fmtspec: fmtspec}
}

func (this *pathnameFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if this.fmtspec == "%s" {
		return append(buf, record.Pathname...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, record.Pathname)...)
	}
}
