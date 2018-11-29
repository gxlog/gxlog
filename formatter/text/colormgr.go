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
	DefaultTraceColor = Green
	DefaultDebugColor = Cyan
	DefaultInfoColor  = Blue
	DefaultWarnColor  = Yellow
	DefaultErrorColor = Red
	DefaultFatalColor = Purple
)

type colorMgr struct {
	colors       []ColorID
	colorEscSeqs [][]byte
	resetEscSeq  []byte
}

func newColorMgr() *colorMgr {
	mgr := &colorMgr{
		colors: []ColorID{
			gxlog.LevelTrace: DefaultTraceColor,
			gxlog.LevelDebug: DefaultDebugColor,
			gxlog.LevelInfo:  DefaultInfoColor,
			gxlog.LevelWarn:  DefaultWarnColor,
			gxlog.LevelError: DefaultErrorColor,
			gxlog.LevelFatal: DefaultFatalColor,
		},
		resetEscSeq: []byte(fmt.Sprintf(cEscSeq, cReset)),
	}
	mgr.initColorSeqs()
	return mgr
}

func (this *colorMgr) initColorSeqs() {
	this.colorEscSeqs = make([][]byte, len(this.colors))
	for i := range this.colors {
		this.colorEscSeqs[i] = makeSeq(this.colors[i])
	}
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

func (this *colorMgr) GetColorEars(level gxlog.LogLevel) ([]byte, []byte) {
	return this.colorEscSeqs[level], this.resetEscSeq
}

func makeSeq(color ColorID) []byte {
	return []byte(fmt.Sprintf(cEscSeq, color))
}
