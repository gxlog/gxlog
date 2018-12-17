package gxlog

type Config struct {
	Level      Level
	TrackLevel Level
	ExitLevel  Level
	Filter     Filter

	Prefix  bool
	Context bool
	Mark    bool
	Limit   bool
}

func NewConfig() *Config {
	return &Config{
		Level:      LevelTrace,
		TrackLevel: LevelFatal,
		ExitLevel:  LevelOff,
		Prefix:     true,
		Context:    true,
		Mark:       true,
		Limit:      true,
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

func (this *Config) WithPrefix(ok bool) *Config {
	this.Prefix = ok
	return this
}

func (this *Config) WithContext(ok bool) *Config {
	this.Context = ok
	return this
}

func (this *Config) WithMark(ok bool) *Config {
	this.Mark = ok
	return this
}

func (this *Config) WithLimit(ok bool) *Config {
	this.Limit = ok
	return this
}
