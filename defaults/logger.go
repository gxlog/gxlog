package defaults

import (
	"os"

	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/writer"
)

var gLogger *gxlog.Logger
var gFormatter *text.Formatter

func init() {
	gLogger = gxlog.New(gxlog.NewConfig())
	gFormatter = text.New(text.NewConfig())
	gLogger.MustLink(gxlog.Slot0, gFormatter, writer.Wrap(os.Stderr))
}

func Logger() *gxlog.Logger {
	return gLogger
}

func Formatter() *text.Formatter {
	return gFormatter
}
