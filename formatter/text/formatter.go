package text

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/gratonos/gxlog"
)

var headerRegexp = regexp.MustCompile("{{([^:%]+?)(?::([^%]*?))?(%.*?)?}}")

type Formatter struct {
	headerFormatters []headerFormatter
}

func New(header string) *Formatter {
	formatter := &Formatter{}
	formatter.SetHeader(header)
	return formatter
}

func (this *Formatter) SetHeader(header string) {
	this.headerFormatters = this.headerFormatters[:0]
	for header != "" {
		indexes := headerRegexp.FindStringSubmatchIndex(header)
		if indexes == nil {
			this.addStaticFormatter(header)
			return
		}
		begin, end := indexes[0], indexes[1]
		if begin != 0 {
			this.addStaticFormatter(header[:begin])
		}
		element, property, fmtspec, ok := extractElement(indexes[2:], header)
		if ok {
			this.addFormatter(element, property, fmtspec)
		}
		header = header[end:]
	}
}

func (this *Formatter) Format(record *gxlog.Record) []byte {
	var buf bytes.Buffer
	for _, f := range this.headerFormatters {
		buf.Write(f.formatHeader(record))
	}
	return buf.Bytes()
}

func (this *Formatter) addStaticFormatter(content string) {
	f := &staticFormatter{content: []byte(content)}
	this.headerFormatters = append(this.headerFormatters, f)
}

func (this *Formatter) addFormatter(element, property, fmtspec string) {
	var f headerFormatter
	switch element {
	case "time":
		f = createTimeFormatter(property, fmtspec)
	case "level":
		f = createLevelFormatter(property, fmtspec)
	case "pathname":
		f = createPathnameFormatter(property, fmtspec)
	case "line":
		f = createLineFormatter(property, fmtspec)
	case "func":
		f = createFuncFormatter(property, fmtspec)
	case "msg":
		f = createMsgFormatter(property, fmtspec)
	}
	if f != nil {
		this.headerFormatters = append(this.headerFormatters, f)
	}
}

func extractElement(indexes []int, header string) (element, property, fmtspec string, ok bool) {
	begin, end := indexes[0], indexes[1]
	if begin < end {
		element = strings.TrimSpace(header[begin:end])
	}
	begin, end = indexes[2], indexes[3]
	if begin < end {
		property = strings.TrimSpace(header[begin:end])
	}
	begin, end = indexes[4], indexes[5]
	if begin < end {
		fmtspec = strings.TrimSpace(header[begin:end])
	}
	if fmtspec == "%" {
		fmtspec = ""
	}
	ok = (element != "")
	return
}
