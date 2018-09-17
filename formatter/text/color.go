package text

import (
	"fmt"

	"github.com/gratonos/gxlog"
)

const (
	escSeq = "\033[%dm"
	reset  = 0
)

type Color int

const (
	Black Color = iota + 30
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

type levelColors struct {
	colors    []Color
	colorSeqs [][]byte
	resetSeq  []byte
}

func newLevelColors() *levelColors {
	colors := &levelColors{
		colors: []Color{
			gxlog.LevelTrace: DefaultTraceColor,
			gxlog.LevelDebug: DefaultDebugColor,
			gxlog.LevelInfo:  DefaultInfoColor,
			gxlog.LevelWarn:  DefaultWarnColor,
			gxlog.LevelError: DefaultErrorColor,
			gxlog.LevelFatal: DefaultFatalColor,
		},
		resetSeq: []byte(fmt.Sprintf(escSeq, reset)),
	}
	colors.initColorSeqs()
	return colors
}

func (this *levelColors) initColorSeqs() {
	this.colorSeqs = make([][]byte, len(this.colors))
	for i := range this.colors {
		this.colorSeqs[i] = makeSeq(this.colors[i])
	}
}

func (this *levelColors) getColor(level gxlog.LogLevel) Color {
	return this.colors[level]
}

func (this *levelColors) setColor(level gxlog.LogLevel, color Color) {
	this.colors[level] = color
	this.colorSeqs[level] = makeSeq(color)
}

func (this *levelColors) updateColors(colors map[gxlog.LogLevel]Color) {
	for level, color := range colors {
		this.setColor(level, color)
	}
}

func (this *levelColors) getColorEars(level gxlog.LogLevel) ([]byte, []byte) {
	return this.colorSeqs[level], this.resetSeq
}

func makeSeq(color Color) []byte {
	return []byte(fmt.Sprintf(escSeq, color))
}
