package gxlog

const (
	DefaultLevel       = LevelTrace
	DefaultExitOnFatal = false
)

type Filter func(*Record) bool

type Config struct {
	Level       Level
	Filter      Filter
	ExitOnFatal bool
}

func NewConfig() *Config {
	return &Config{
		Level:       DefaultLevel,
		ExitOnFatal: DefaultExitOnFatal,
	}
}

func (this *Config) WithLevel(level Level) *Config {
	this.Level = level
	return this
}

func (this *Config) WithFilter(filter Filter) *Config {
	this.Filter = filter
	return this
}

func (this *Config) WithExitOnFatal(ok bool) *Config {
	this.ExitOnFatal = ok
	return this
}
