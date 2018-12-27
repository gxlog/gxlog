package text

import (
	"fmt"
	"strconv"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/internal/util"
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

func (formatter *pkgFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	pkg := util.LastSegments(record.Pkg, formatter.segments, '/')
	if formatter.fmtspec == "%s" {
		return append(buf, pkg...)
	} else {
		return append(buf, fmt.Sprintf(formatter.fmtspec, pkg)...)
	}
}
