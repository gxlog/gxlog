package gxlog

const (
	DefaultLevel       = LevelTrace
	DefaultExitOnFatal = false
)

type Config struct {
	Level       LogLevel
	ExitOnFatal bool
}

func NewConfig() *Config {
	return &Config{
		Level:       DefaultLevel,
		ExitOnFatal: DefaultExitOnFatal,
	}
}

func (this *Config) WithLevel(level LogLevel) *Config {
	this.Level = level
	return this
}

func (this *Config) WithExitOnFatal(ok bool) *Config {
	this.ExitOnFatal = ok
	return this
}
