package text

import (
	"github.com/gratonos/gxlog"
)

const (
	DefaultHeader = "{{time}} {{level}} {{file}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	DefaultMinBufSize  = 256
	DefaultEnableColor = false
)

type Config struct {
	Header      string
	MinBufSize  int
	ColorMap    map[gxlog.LogLevel]ColorID
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

func (this *Config) WithColorMap(colorMap map[gxlog.LogLevel]ColorID) *Config {
	this.ColorMap = colorMap
	return this
}

func (this *Config) WithEnableColor(enable bool) *Config {
	this.EnableColor = enable
	return this
}
