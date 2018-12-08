package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

type ColorID int

const (
	Black ColorID = iota + 30
	Red
	Green
	Yellow
	Blue
	Purple
	Cyan
	White
)

const (
	DefaultTraceColor  = Green
	DefaultDebugColor  = Green
	DefaultInfoColor   = Green
	DefaultWarnColor   = Yellow
	DefaultErrorColor  = Red
	DefaultFatalColor  = Red
	DefaultMarkedColor = Purple
)

const (
	cEscSeq = "\033[%dm"
	cReset  = 0
)

type colorMgr struct {
	colors      []ColorID
	markedColor ColorID

	colorSeqs [][]byte
	markedSeq []byte
	resetSeq  []byte
}

func newColorMgr() *colorMgr {
	colors := []ColorID{
		gxlog.LevelTrace: DefaultTraceColor,
		gxlog.LevelDebug: DefaultDebugColor,
		gxlog.LevelInfo:  DefaultInfoColor,
		gxlog.LevelWarn:  DefaultWarnColor,
		gxlog.LevelError: DefaultErrorColor,
		gxlog.LevelFatal: DefaultFatalColor,
	}
	mgr := &colorMgr{
		colors:      colors,
		markedColor: DefaultMarkedColor,
		colorSeqs:   initColorSeqs(colors),
		markedSeq:   makeSeq(DefaultMarkedColor),
		resetSeq:    makeSeq(0),
	}
	return mgr
}

func (this *colorMgr) GetColor(level gxlog.Level) ColorID {
	return this.colors[level]
}

func (this *colorMgr) SetColor(level gxlog.Level, color ColorID) {
	this.colors[level] = color
	this.colorSeqs[level] = makeSeq(color)
}

func (this *colorMgr) MapColors(colorMap map[gxlog.Level]ColorID) {
	for level, color := range colorMap {
		this.SetColor(level, color)
	}
}

func (this *colorMgr) GetMarkedColor() ColorID {
	return this.markedColor
}

func (this *colorMgr) SetMarkedColor(color ColorID) {
	this.markedColor = color
	this.markedSeq = makeSeq(color)
}

func (this *colorMgr) GetColorEars(level gxlog.Level) ([]byte, []byte) {
	return this.colorSeqs[level], this.resetSeq
}

func (this *colorMgr) GetMarkedColorEars() ([]byte, []byte) {
	return this.markedSeq, this.resetSeq
}

func initColorSeqs(colors []ColorID) [][]byte {
	colorSeqs := make([][]byte, len(colors))
	for i := range colors {
		colorSeqs[i] = makeSeq(colors[i])
	}
	return colorSeqs
}

func makeSeq(color ColorID) []byte {
	return []byte(fmt.Sprintf(cEscSeq, color))
}
