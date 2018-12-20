package text

import (
	"regexp"
	"strings"
	"sync"

	"github.com/gxlog/gxlog"
)

var gHeaderRegexp = regexp.MustCompile("{{([^:%]*?)(?::([^%]*?))?(%.*?)?}}")

type Formatter struct {
	header      string
	minBufSize  int
	enableColor bool

	colorMgr  *colorMgr
	appenders []*headerAppender
	suffix    string

	lock sync.Mutex
}

func New(config *Config) *Formatter {
	if config.MinBufSize < 0 {
		panic("formatter/text.New: Config.MinBufSize must not be negative")
	}
	formatter := &Formatter{
		minBufSize:  config.MinBufSize,
		enableColor: config.EnableColor,
		colorMgr:    newColorMgr(),
	}
	formatter.SetHeader(config.Header)
	formatter.MapColors(config.ColorMap)
	return formatter
}

func (this *Formatter) Header() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header
}

func (this *Formatter) SetHeader(header string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.header = header
	this.appenders = this.appenders[:0]
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
	this.suffix = staticText + header
}

func (this *Formatter) MinBufSize() int {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.minBufSize
}

func (this *Formatter) SetMinBufSize(size int) {
	if size < 0 {
		panic("formatter/text.SetMinBufSize: size must not be negative")
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	this.minBufSize = size
}

func (this *Formatter) ColorEnabled() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.enableColor
}

func (this *Formatter) EnableColor() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.enableColor = true
}

func (this *Formatter) DisableColor() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.enableColor = false
}

func (this *Formatter) Color(level gxlog.Level) ColorID {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.colorMgr.Color(level)
}

func (this *Formatter) SetColor(level gxlog.Level, color ColorID) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.colorMgr.SetColor(level, color)
}

func (this *Formatter) MapColors(colorMap map[gxlog.Level]ColorID) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.colorMgr.MapColors(colorMap)
}

func (this *Formatter) MarkedColor() ColorID {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.colorMgr.MarkedColor()
}

func (this *Formatter) SetMarkedColor(color ColorID) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.colorMgr.SetMarkedColor(color)
}

func (this *Formatter) Format(record *gxlog.Record) []byte {
	this.lock.Lock()
	defer this.lock.Unlock()

	var left, right []byte
	if this.enableColor {
		if record.Aux.Marked {
			left, right = this.colorMgr.MarkedColorEars()
		} else {
			left, right = this.colorMgr.ColorEars(record.Level)
		}
	}

	buf := make([]byte, 0, this.minBufSize)
	buf = append(buf, left...)
	for _, appender := range this.appenders {
		buf = appender.AppendHeader(buf, record)
	}
	buf = append(buf, this.suffix...)
	buf = append(buf, right...)

	return buf
}

func (this *Formatter) addAppender(element, property, fmtspec, staticText string) bool {
	appender := newHeaderAppender(element, property, fmtspec, staticText)
	if appender == nil {
		return false
	}
	this.appenders = append(this.appenders, appender)
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
