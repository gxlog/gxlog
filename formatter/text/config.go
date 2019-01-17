package text

import (
	"github.com/gxlog/gxlog/iface"
)

// All predefined headers here.
const (
	FullHeader = "{{time}} {{level}} {{file}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	CompactHeader = "{{time:time.us}} {{level}} {{file:1}}:{{line}} " +
		"{{pkg}}.{{func}} {{prefix}}[{{context}}] {{msg}}\n"
	SyslogHeader = "{{file:1}}:{{line}} {{pkg}}.{{func}} {{prefix}}" +
		"[{{context}}] {{msg}}\n"
)

// A Config is used to configure a text formatter.
type Config struct {
	// Header is the format specifier of a text formatter.
	// It is used to specify which and how the fields of a Record to be formatted.
	// The pattern of a field specifier is {{<name>[:property][%fmtstr]}}.
	// e.g. {{level:char}}, {{line%05d}}, {{pkg:1}}, {{context:list%40s}}
	// All fields have support for the fmtstr. If the fmtstr is NOT the default
	// one of a field, it will be passed to fmt.Sprintf to format the field and
	// this affects the performance a little.
	// The supported properties vary with fields.
	// All supported fields are as the follows:
	//   name    | supported property       | defaults     | property examples
	// ----------+--------------------------+--------------+------------------------
	//   time    | <date|time>[.ms|.us|.ns] | "date.us" %s | "date.ns", "time"
	//           | layout that is supported |              | time.RFC3339Nano
	//           |   by the time package    |              | "02 Jan 06 15:04 -0700"
	//   level   | <full|char>              | "full"    %s | "full", "char"
	//   file    | <lastSegs>               | 0         %s | 0, 1, 2, ...
	//   line    |                          |           %d |
	//   pkg     | <lastSegs>               | 0         %s | 0, 1, 2, ...
	//   func    | <lastSegs>               | 0         %s | 0, 1, 2, ...
	//   prefix  |                          |           %s |
	//   context | <pair|list>              | "pair"    %s | "pair", "list"
	//   msg     |                          |           %s |
	// If Header is not specified, FullHeader is used.
	Header string
	// MinBufSize is the initial size of the internal buf of a formatter.
	// MinBufSize must NOT be negative. If it is not specified, 256 is used.
	MinBufSize int
	// ColorMap is used to remap the color of each level.
	// By default, the color of Trace, Debug and Info is Green, the color of Warn
	// is Yellow, the color of Error and Fatal is Red. The color of a marked log
	// is Magenta despite of its level.
	// The color of a level is left to be unchanged if it is not in the map.
	ColorMap map[iface.Level]Color
	// EnableColor enables colorization if it is true.
	EnableColor bool
}

func (config *Config) setDefaults() {
	if config.Header == "" {
		config.Header = FullHeader
	}
	if config.MinBufSize == 0 {
		config.MinBufSize = 256
	}
}
