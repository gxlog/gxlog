package text

import (
	"github.com/gratonos/gxlog"
)

const (
	DefaultHeader = "{{time}} {{level}} {{pathname}}:{{line}} {{func}} " +
		"{{prefix}}{{context}} {{msg}}\n"
	DefaultMinBufSize    = 256
	DefaultBatchBufCount = 16
	DefaultEnableColor   = false
)

type Config struct {
	Header        string
	MinBufSize    int
	BatchBufCount int
	ColorMap      map[gxlog.LogLevel]ColorID
	EnableColor   bool
}

func NewConfig() *Config {
	return &Config{
		Header:        DefaultHeader,
		MinBufSize:    DefaultMinBufSize,
		BatchBufCount: DefaultBatchBufCount,
		EnableColor:   DefaultEnableColor,
	}
}

func (this *Config) WithHeader(header string) *Config {
	this.Header = header
	return this
}

func (this *Config) WithBufConfig(minBufSize, batchBufCount int) *Config {
	this.MinBufSize = minBufSize
	this.BatchBufCount = batchBufCount
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
