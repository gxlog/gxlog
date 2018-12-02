package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type msgFormatter struct {
	property string
	fmtspec  string
}

func newMsgFormatter(property, fmtspec string) *msgFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &msgFormatter{property: property, fmtspec: fmtspec}
}

func (this *msgFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	if this.fmtspec == "%s" {
		return append(buf, record.Msg...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, record.Msg)...)
	}
}
