package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

const (
	cEscSeq = "\033[%dm"
	cReset  = 0
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
	DefaultFatalColor  = Purple
	DefaultMarkedColor = Blue
)

type colorMgr struct {
	colors      []ColorID
	markedColor ColorID

	colorEscSeqs [][]byte
	markedEscSeq []byte
	resetEscSeq  []byte
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
		colors:       colors,
		markedColor:  DefaultMarkedColor,
		colorEscSeqs: initColorSeqs(colors),
		markedEscSeq: makeSeq(DefaultMarkedColor),
		resetEscSeq:  makeSeq(0),
	}
	return mgr
}

func (this *colorMgr) GetColor(level gxlog.LogLevel) ColorID {
	return this.colors[level]
}

func (this *colorMgr) SetColor(level gxlog.LogLevel, color ColorID) {
	this.colors[level] = color
	this.colorEscSeqs[level] = makeSeq(color)
}

func (this *colorMgr) MapColors(colorMap map[gxlog.LogLevel]ColorID) {
	for level, color := range colorMap {
		this.SetColor(level, color)
	}
}

func (this *colorMgr) GetMarkedColor() ColorID {
	return this.markedColor
}

func (this *colorMgr) SetMarkedColor(color ColorID) {
	this.markedColor = color
	this.markedEscSeq = makeSeq(color)
}

func (this *colorMgr) GetColorEars(level gxlog.LogLevel) ([]byte, []byte) {
	return this.colorEscSeqs[level], this.resetEscSeq
}

func (this *colorMgr) GetMarkedColorEars() ([]byte, []byte) {
	return this.markedEscSeq, this.resetEscSeq
}

func initColorSeqs(colors []ColorID) [][]byte {
	colorEscSeqs := make([][]byte, len(colors))
	for i := range colors {
		colorEscSeqs[i] = makeSeq(colors[i])
	}
	return colorEscSeqs
}

func makeSeq(color ColorID) []byte {
	return []byte(fmt.Sprintf(cEscSeq, color))
}
