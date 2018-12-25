package gxlog

type Flag int

const (
	Prefix Flag = 0x01 << iota
	Contexts
	DynamicContexts
	Mark
	Limit
)

type Filter func(*Record) bool

type Config struct {
	Level      Level
	TimeLevel  Level
	PanicLevel Level
	TrackLevel Level
	ExitLevel  Level
	Filter     Filter
	Flags      Flag
}

func NewConfig() *Config {
	return &Config{
		Level:      LevelTrace,
		TimeLevel:  LevelTrace,
		PanicLevel: LevelFatal,
		TrackLevel: LevelFatal,
		ExitLevel:  LevelOff,
		Flags:      Prefix | Contexts | DynamicContexts | Mark | Limit,
	}
}

func (this *Config) WithLevel(level Level) *Config {
	this.Level = level
	return this
}

func (this *Config) WithTimeLevel(level Level) *Config {
	this.TimeLevel = level
	return this
}

func (this *Config) WithPanicLevel(level Level) *Config {
	this.PanicLevel = level
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

func (this *Config) WithFlags(flags Flag) *Config {
	this.Flags = flags
	return this
}

func (this *Config) WithEnabled(flags Flag) *Config {
	this.Flags |= flags
	return this
}

func (this *Config) WithDisabled(flags Flag) *Config {
	this.Flags &^= flags
	return this
}
