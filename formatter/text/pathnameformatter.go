package text

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gratonos/gxlog"
)

type pathnameFormatter struct {
	segments int
	fmtspec  string
}

func newPathnameFormatter(property, fmtspec string) *pathnameFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &pathnameFormatter{
		segments: segments,
		fmtspec:  fmtspec,
	}
}

func (this *pathnameFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	pathname := filepath.ToSlash(record.Pathname)
	if this.segments > 0 {
		pathname = lastSegments(pathname, this.segments, '/')
	}
	if this.fmtspec == "%s" {
		return append(buf, pathname...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, pathname)...)
	}
}
