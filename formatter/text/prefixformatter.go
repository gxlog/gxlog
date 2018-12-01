package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
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

func (this *prefixFormatter) FormatElement(record *gxlog.Record) string {
	if this.fmtspec == "%s" {
		return record.Prefix
	} else {
		return fmt.Sprintf(this.fmtspec, record.Prefix)
	}
}
