// Package defaults provides the default Logger and the default Formatter.
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

// Logger returns the default Logger. The default Logger has the default
// Formatter (a text formatter) and a writer wrapper of os.Stderr linked
// in Slot0. The rest slots are free.
func Logger() *gxlog.Logger {
	return defaultLogger
}

// Formatter returns the default Formatter. It is a text formatter.
func Formatter() *text.Formatter {
	return defaultFormatter
}
