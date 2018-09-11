package text

import (
	"regexp"
	"strings"

	"github.com/gratonos/gxlog"
)

var headerRegexp = regexp.MustCompile("{{([^:%]*?)(?::([^%]*?))?(%.*?)?}}")

type Formatter struct {
	headerAppenders []headerAppender
	buf             []byte
}

func New(header string) *Formatter {
	formatter := &Formatter{}
	formatter.SetHeader(header)
	return formatter
}

func (this *Formatter) SetHeader(header string) {
	this.headerAppenders = this.headerAppenders[:0]
	for header != "" {
		indexes := headerRegexp.FindStringSubmatchIndex(header)
		if indexes == nil {
			this.addStaticAppender(header)
			return
		}
		begin, end := indexes[0], indexes[1]
		if begin != 0 {
			this.addStaticAppender(header[:begin])
		}
		element, property, fmtspec, ok := extractElement(indexes[2:], header)
		if ok {
			this.addAppender(element, property, fmtspec)
		}
		header = header[end:]
	}
}

func (this *Formatter) Format(record *gxlog.Record) []byte {
	this.buf = this.buf[:0]
	for _, f := range this.headerAppenders {
		this.buf = f.appendHeader(this.buf, record)
	}
	return this.buf
}

func (this *Formatter) addStaticAppender(content string) {
	f := &staticAppender{content: []byte(content)}
	this.headerAppenders = append(this.headerAppenders, f)
}

func (this *Formatter) addAppender(element, property, fmtspec string) {
	var f headerAppender
	switch element {
	case "time":
		f = createTimeAppender(property, fmtspec)
	case "level":
		f = createLevelAppender(property, fmtspec)
	case "pathname":
		f = createPathnameAppender(property, fmtspec)
	case "line":
		f = createLineAppender(property, fmtspec)
	case "func":
		f = createFuncAppender(property, fmtspec)
	case "msg":
		f = createMsgAppender(property, fmtspec)
	}
	if f != nil {
		this.headerAppenders = append(this.headerAppenders, f)
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
