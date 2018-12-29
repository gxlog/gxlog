package text

import (
	"fmt"

	"github.com/gxlog/gxlog"
)

// The Color defines the color type.
type Color int

// All available colors here.
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const escSeqFmt = "\033[%dm"

type colorMgr struct {
	colors      []Color
	markedColor Color

	colorSeqs [][]byte
	markedSeq []byte
	resetSeq  []byte
}

func newColorMgr() *colorMgr {
	colors := []Color{
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

func (mgr *colorMgr) Color(level gxlog.Level) Color {
	return mgr.colors[level]
}

func (mgr *colorMgr) SetColor(level gxlog.Level, color Color) {
	mgr.colors[level] = color
	mgr.colorSeqs[level] = makeSeq(color)
}

func (mgr *colorMgr) MapColors(colorMap map[gxlog.Level]Color) {
	for level, color := range colorMap {
		mgr.SetColor(level, color)
	}
}

func (mgr *colorMgr) MarkedColor() Color {
	return mgr.markedColor
}

func (mgr *colorMgr) SetMarkedColor(color Color) {
	mgr.markedColor = color
	mgr.markedSeq = makeSeq(color)
}

func (mgr *colorMgr) ColorEars(level gxlog.Level) ([]byte, []byte) {
	return mgr.colorSeqs[level], mgr.resetSeq
}

func (mgr *colorMgr) MarkedColorEars() ([]byte, []byte) {
	return mgr.markedSeq, mgr.resetSeq
}

func initColorSeqs(colors []Color) [][]byte {
	colorSeqs := make([][]byte, len(colors))
	for i := range colors {
		colorSeqs[i] = makeSeq(colors[i])
	}
	return colorSeqs
}

func makeSeq(color Color) []byte {
	return []byte(fmt.Sprintf(escSeqFmt, color))
}
