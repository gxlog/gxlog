package defaults

import (
	"os"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/text"
	"github.com/gxlog/gxlog/writer"
)

var gLogger *gxlog.Logger
var gFormatter *text.Formatter

func init() {
	gLogger = gxlog.New(gxlog.NewConfig())
	gFormatter = text.New(text.NewConfig())
	gLogger.Link(gxlog.Slot0, gFormatter, writer.Wrap(os.Stderr))
}

func Logger() *gxlog.Logger {
	return gLogger
}

func Formatter() *text.Formatter {
	return gFormatter
}
