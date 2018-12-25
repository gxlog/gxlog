package text

import (
	"fmt"

	"github.com/gxlog/gxlog"
)

type ColorID int

const (
	Black ColorID = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
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
		gxlog.Trace: Green,
		gxlog.Debug: Green,
		gxlog.Info:  Green,
		gxlog.Warn:  Yellow,
		gxlog.Error: Red,
		gxlog.Fatal: Red,
	}
	mgr := &colorMgr{
		colors:      colors,
		markedColor: Magenta,
		colorSeqs:   initColorSeqs(colors),
		markedSeq:   makeSeq(Magenta),
		resetSeq:    makeSeq(0),
	}
	return mgr
}

func (this *colorMgr) Color(level gxlog.Level) ColorID {
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

func (this *colorMgr) MarkedColor() ColorID {
	return this.markedColor
}

func (this *colorMgr) SetMarkedColor(color ColorID) {
	this.markedColor = color
	this.markedSeq = makeSeq(color)
}

func (this *colorMgr) ColorEars(level gxlog.Level) ([]byte, []byte) {
	return this.colorSeqs[level], this.resetSeq
}

func (this *colorMgr) MarkedColorEars() ([]byte, []byte) {
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
