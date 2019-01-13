package text

import (
	"fmt"
	"strconv"

	"github.com/gxlog/gxlog/formatter/internal/util"
	"github.com/gxlog/gxlog/iface"
)

type pkgFormatter struct {
	segments int
	fmtspec  string
}

func newPkgFormatter(property, fmtspec string) *pkgFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &pkgFormatter{
		segments: segments,
		fmtspec:  fmtspec,
	}
}

func (formatter *pkgFormatter) FormatElement(buf []byte,
	record *iface.Record) []byte {
	pkg := util.LastSegments(record.Pkg, formatter.segments, '/')
	if formatter.fmtspec == "%s" {
		return append(buf, pkg...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, pkg)...)
}
