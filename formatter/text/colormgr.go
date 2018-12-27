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

const escSeqFmt = "\033[%dm"

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

func (mgr *colorMgr) Color(level gxlog.Level) ColorID {
	return mgr.colors[level]
}

func (mgr *colorMgr) SetColor(level gxlog.Level, color ColorID) {
	mgr.colors[level] = color
	mgr.colorSeqs[level] = makeSeq(color)
}

func (mgr *colorMgr) MapColors(colorMap map[gxlog.Level]ColorID) {
	for level, color := range colorMap {
		mgr.SetColor(level, color)
	}
}

func (mgr *colorMgr) MarkedColor() ColorID {
	return mgr.markedColor
}

func (mgr *colorMgr) SetMarkedColor(color ColorID) {
	mgr.markedColor = color
	mgr.markedSeq = makeSeq(color)
}

func (mgr *colorMgr) ColorEars(level gxlog.Level) ([]byte, []byte) {
	return mgr.colorSeqs[level], mgr.resetSeq
}

func (mgr *colorMgr) MarkedColorEars() ([]byte, []byte) {
	return mgr.markedSeq, mgr.resetSeq
}

func initColorSeqs(colors []ColorID) [][]byte {
	colorSeqs := make([][]byte, len(colors))
	for i := range colors {
		colorSeqs[i] = makeSeq(colors[i])
	}
	return colorSeqs
}

func makeSeq(color ColorID) []byte {
	return []byte(fmt.Sprintf(escSeqFmt, color))
}
