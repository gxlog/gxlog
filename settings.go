package gxlog

func (this *Logger) Config() *Config {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := *this.config
	return &copyConfig
}

func (this *Logger) SetConfig(config *Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := *config
	this.config = &copyConfig
	return nil
}

func (this *Logger) UpdateConfig(fn func(Config) Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	config := fn(*this.config)
	this.config = &config
	return nil
}

func (this *Logger) Level() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level
}

func (this *Logger) SetLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Level = level
}

func (this *Logger) TimeLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.TimeLevel
}

func (this *Logger) SetTimeLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.TimeLevel = level
}

func (this *Logger) PanicLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.PanicLevel
}

func (this *Logger) SetPanicLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.PanicLevel = level
}

func (this *Logger) TrackLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.TrackLevel
}

func (this *Logger) SetTrackLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.TrackLevel = level
}

func (this *Logger) ExitLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.ExitLevel
}

func (this *Logger) SetExitLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.ExitLevel = level
}

func (this *Logger) Filter() Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Filter
}

func (this *Logger) SetFilter(filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Filter = filter
}

func (this *Logger) Prefix() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Prefix
}

func (this *Logger) SetPrefix(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Prefix = ok
}

func (this *Logger) Context() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Context
}

func (this *Logger) SetContext(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Context = ok
}

func (this *Logger) Dynamic() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Dynamic
}

func (this *Logger) SetDynamic(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Dynamic = ok
}

func (this *Logger) Mark() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Mark
}

func (this *Logger) SetMark(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Mark = ok
}

func (this *Logger) Limit() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Limit
}

func (this *Logger) SetLimit(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Limit = ok
}
