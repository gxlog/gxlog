package text

import (
	"github.com/gxlog/gxlog"
)

// All predefined headers here.
const (
	DefaultHeader = "{{time}} {{level}} {{file}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	CompactHeader = "{{time:time.us}} {{level}} {{file:1}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	SyslogHeader = "{{file:1}}:{{line}} {{pkg}}.{{func}} {{prefix}}[{{context}}] {{msg}}\n"
)

// A Config is used to configure a text formatter.
// A Config should be created with NewConfig.
type Config struct {
	// Header is the format specifier of a text formatter.
	// It is used to specify which and how the fields of a Record to be formatted.
	// The pattern of a field specifier is {{<name>[:property][%fmtstr]}}.
	// e.g. {{level:char}}, {{line%05d}}, {{pkg:1}}, {{context:list%40s}}
	// All fields have support for the fmtstr. The fmtstr will be passed to
	// fmt.Sprintf to format the field.
	// The supported properties vary with fields.
	// All supported fields are as follows:
	//            supported property          defaults      property examples
	// --------------------------------------------------------------------------
	//   time     <date|time>[.ms|.us|.ns]    date.us %s    date.ns, time
	//            layout that is supported                  time.RFC3339Nano
	//              by the time package                     02 Jan 06 15:04 -0700
	//   level    <full|char>                 full    %s    full, char
	//   file     <lastSegs>                  0       %s    0, 1, 2, ...
	//   line                                         %d
	//   pkg      <lastSegs>                  0       %s    0, 1, 2, ...
	//   func     <lastSegs>                  0       %s    0, 1, 2, ...
	//   prefix                                       %s
	//   context  <pair|list>                 pair    %s    pair, list
	//   msg                                          %s
	Header string
	// MinBufSize is the initial size of the internal buf of a formatter.
	// MinBufSize must not be negative.
	MinBufSize int
	// ColorMap is used to remap the color of each level.
	// The color of a level is left to be unchanged if it is not in the map.
	// By default, the color of level Trace, Debug and Info is Green, the color
	// of level Warn is Yellow, the color of level Error and Fatal is Red and
	// the color of a marked log is Magenta no matter at which level it is.
	ColorMap map[gxlog.Level]ColorID
	// EnableColor enables colorization if it is true.
	EnableColor bool
}

// NewConfig creates a new Config. By default, the Header is the DefaultHeader,
// the MinBufSize is 256, the ColorMap is nil, the EnableColor is false.
func NewConfig() *Config {
	return &Config{
		Header:     DefaultHeader,
		MinBufSize: 256,
	}
}

// WithHeader sets the Header of the Config and returns it.
func (cfg *Config) WithHeader(header string) *Config {
	cfg.Header = header
	return cfg
}

// WithMinBufSize sets the MinBufSize of the Config and returns it.
func (cfg *Config) WithMinBufSize(size int) *Config {
	cfg.MinBufSize = size
	return cfg
}

// WithColorMap sets the ColorMap of the Config and returns it.
func (cfg *Config) WithColorMap(colorMap map[gxlog.Level]ColorID) *Config {
	cfg.ColorMap = colorMap
	return cfg
}

// WithEnableColor sets the EnableColor of the Config and returns it.
func (cfg *Config) WithEnableColor(enable bool) *Config {
	cfg.EnableColor = enable
	return cfg
}
