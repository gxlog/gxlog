package text

import (
	"github.com/gxlog/gxlog"
)

const (
	DefaultHeader = "{{time}} {{level}} {{file}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	CompactHeader = "{{time:time.us}} {{level}} {{file:1}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	SyslogHeader = "{{file:1}}:{{line}} {{pkg}}.{{func}} {{prefix}}[{{context}}] {{msg}}\n"
)

type Config struct {
	Header      string
	MinBufSize  int
	ColorMap    map[gxlog.Level]ColorID
	EnableColor bool
}

func NewConfig() *Config {
	return &Config{
		Header:     DefaultHeader,
		MinBufSize: 256,
	}
}

func (cfg *Config) WithHeader(header string) *Config {
	cfg.Header = header
	return cfg
}

func (cfg *Config) WithMinBufSize(size int) *Config {
	cfg.MinBufSize = size
	return cfg
}

func (cfg *Config) WithColorMap(colorMap map[gxlog.Level]ColorID) *Config {
	cfg.ColorMap = colorMap
	return cfg
}

func (cfg *Config) WithEnableColor(enable bool) *Config {
	cfg.EnableColor = enable
	return cfg
}
