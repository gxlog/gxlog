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

func (this *pathnameFormatter) FormatElement(record *gxlog.Record) string {
	if this.fmtspec == "%s" {
		return record.Pathname
	} else {
		return fmt.Sprintf(this.fmtspec, record.Pathname)
	}
}
