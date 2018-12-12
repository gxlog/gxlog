package text

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gxlog/gxlog"
)

type fileFormatter struct {
	segments int
	fmtspec  string
}

func newFileFormatter(property, fmtspec string) *fileFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &fileFormatter{
		segments: segments,
		fmtspec:  fmtspec,
	}
}

func (this *fileFormatter) FormatElement(buf []byte, record *gxlog.Record) []byte {
	file := filepath.ToSlash(record.File)
	if this.segments > 0 {
		file = lastSegments(file, this.segments, '/')
	}
	if this.fmtspec == "%s" {
		return append(buf, file...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, file)...)
	}
}
