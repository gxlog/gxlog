package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog"
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

func (this *pkgFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	pkg := record.Pkg
	if this.segments > 0 {
		pkg = lastSegments(pkg, this.segments, '/')
	}
	if this.fmtspec == "%s" {
		return append(buf, pkg...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, pkg)...)
	}
}
