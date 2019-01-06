package text

import (
	"fmt"

	"github.com/gxlog/gxlog/iface"
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

func (formatter *msgFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if formatter.fmtspec == "%s" {
		return append(buf, record.Msg...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, record.Msg)...)
}
