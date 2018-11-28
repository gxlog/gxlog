package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type funcFormatter struct {
	property string
	fmtspec  string
}

func newFuncFormatter(property, fmtspec string) *funcFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &funcFormatter{property: property, fmtspec: fmtspec}
}

func (this *funcFormatter) FormatElement(record *gxlog.Record) string {
	if this.fmtspec == "%s" {
		return record.Func
	} else {
		return fmt.Sprintf(this.fmtspec, record.Func)
	}
}
