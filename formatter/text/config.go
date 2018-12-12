package text

import (
	"github.com/gxlog/gxlog"
)

const (
	DefaultHeader = "{{time}} {{level}} {{file}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	DefaultMinBufSize  = 256
	DefaultEnableColor = false
)

const CompactHeader = "{{time:time.us}} {{level}} {{file:1}}:{{line}} {{pkg}}.{{func}} " +
	"{{prefix}}[{{context}}] {{msg}}\n"

type Config struct {
	Header      string
	MinBufSize  int
	ColorMap    map[gxlog.Level]ColorID
	EnableColor bool
}

func NewConfig() *Config {
	return &Config{
		Header:      DefaultHeader,
		MinBufSize:  DefaultMinBufSize,
		EnableColor: DefaultEnableColor,
	}
}

func (this *Config) WithHeader(header string) *Config {
	this.Header = header
	return this
}

func (this *Config) WithMinBufSize(size int) *Config {
	this.MinBufSize = size
	return this
}

func (this *Config) WithColorMap(colorMap map[gxlog.Level]ColorID) *Config {
	this.ColorMap = colorMap
	return this
}

func (this *Config) WithEnableColor(enable bool) *Config {
	this.EnableColor = enable
	return this
}
