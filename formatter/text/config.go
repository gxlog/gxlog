package text

import (
	"github.com/gratonos/gxlog"
)

const (
	DefaultHeader      = "{{time}} {{level}} {{pathname}}:{{line}} {{func}} {{prefix}}{{msg}}\n"
	DefaultEnableColor = true
)

type Config struct {
	Header      string
	ColorMap    map[gxlog.LogLevel]ColorID
	EnableColor bool
}

func NewConfig() *Config {
	return &Config{
		Header:      DefaultHeader,
		EnableColor: DefaultEnableColor,
	}
}

func (this *Config) WithHeader(header string) *Config {
	this.Header = header
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
