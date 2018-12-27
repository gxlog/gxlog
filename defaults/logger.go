package defaults

import (
	"os"

	"github.com/gxlog/gxlog"
	"github.com/gxlog/gxlog/formatter/text"
	"github.com/gxlog/gxlog/writer"
)

var defaultLogger *gxlog.Logger
var defaultFormatter *text.Formatter

func init() {
	defaultLogger = gxlog.New(gxlog.NewConfig())
	defaultFormatter = text.New(text.NewConfig())
	defaultLogger.Link(gxlog.Slot0, defaultFormatter, writer.Wrap(os.Stderr))
}

func Logger() *gxlog.Logger {
	return defaultLogger
}

func Formatter() *text.Formatter {
	return defaultFormatter
}
