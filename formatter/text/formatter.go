package text

import (
	"regexp"
	"strings"

	"github.com/gratonos/gxlog"
)

var gHeaderRegexp = regexp.MustCompile("{{([^:%]*?)(?::([^%]*?))?(%.*?)?}}")

type Formatter struct {
	headerAppenders []*headerAppender
	suffix          []byte
	buf             []byte
	*colorMgr
	enableColor bool
}

func New(config *Config) *Formatter {
	formatter := &Formatter{
		colorMgr: newColorMgr(),
	}
	formatter.SetHeader(config.Header)
	formatter.MapColors(config.ColorMap)
	formatter.enableColor = config.EnableColor
	return formatter
}

func (this *Formatter) SetHeader(header string) {
	this.headerAppenders = this.headerAppenders[:0]
	var staticText string
	for header != "" {
		indexes := gHeaderRegexp.FindStringSubmatchIndex(header)
		if indexes == nil {
			break
		}
		begin, end := indexes[0], indexes[1]
		staticText += header[:begin]
		element, property, fmtspec := extractElement(indexes, header)
		if this.addAppender(element, property, fmtspec, staticText) {
			staticText = ""
		}
		header = header[end:]
	}
	this.suffix = []byte(staticText + header)
}

func (this *Formatter) EnableColor() {
	this.enableColor = true
}

func (this *Formatter) DisableColor() {
	this.enableColor = false
}

func (this *Formatter) Format(record *gxlog.Record) []byte {
	var left, right []byte
	if this.enableColor {
		if record.Marked {
			left, right = this.colorMgr.getMarkedColorEars()
		} else {
			left, right = this.colorMgr.getColorEars(record.Level)
		}
	}
	this.buf = this.buf[:0]
	this.buf = append(this.buf, left...)
	for _, appender := range this.headerAppenders {
		this.buf = appender.AppendHeader(this.buf, record)
	}
	this.buf = append(this.buf, this.suffix...)
	this.buf = append(this.buf, right...)
	return this.buf
}

func (this *Formatter) addAppender(element, property, fmtspec, staticText string) bool {
	appender := newHeaderAppender(element, property, fmtspec, staticText)
	if appender == nil {
		return false
	}
	this.headerAppenders = append(this.headerAppenders, appender)
	return true
}

func extractElement(indexes []int, header string) (element, property, fmtspec string) {
	element = strings.ToLower(getField(header, indexes[2], indexes[3]))
	property = getField(header, indexes[4], indexes[5])
	fmtspec = getField(header, indexes[6], indexes[7])
	if fmtspec == "%" {
		fmtspec = ""
	}
	return element, property, fmtspec
}

func getField(str string, begin, end int) string {
	if begin < end {
		return strings.TrimSpace(str[begin:end])
	}
	return ""
}
