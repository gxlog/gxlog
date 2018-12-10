package gxlog

const (
	DefaultLevel      = LevelTrace
	DefaultTrackLevel = LevelFatal
	DefaultExitLevel  = LevelOff
)

type Filter func(*Record) bool

type Config struct {
	Level      Level
	TrackLevel Level
	ExitLevel  Level
	Filter     Filter
}

func NewConfig() *Config {
	return &Config{
		Level:      DefaultLevel,
		TrackLevel: DefaultTrackLevel,
		ExitLevel:  DefaultExitLevel,
	}
}

func (this *Config) WithLevel(level Level) *Config {
	this.Level = level
	return this
}

func (this *Config) WithTrackLevel(level Level) *Config {
	this.TrackLevel = level
	return this
}
func (this *Config) WithExitLevel(level Level) *Config {
	this.ExitLevel = level
	return this
}
func (this *Config) WithFilter(filter Filter) *Config {
	this.Filter = filter
	return this
}
