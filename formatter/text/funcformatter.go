package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog"
)

type funcFormatter struct {
	segments int
	fmtspec  string
}

func newFuncFormatter(property, fmtspec string) *funcFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &funcFormatter{
		segments: segments,
		fmtspec:  fmtspec,
	}
}

func (this *funcFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	fn := record.Func
	if this.segments > 0 {
		fn = lastSegments(fn, this.segments, '.')
	}
	if this.fmtspec == "%s" {
		return append(buf, fn...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, fn)...)
	}
}
