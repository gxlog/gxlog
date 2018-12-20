package gxlog

func (this *logger) Config() *Config {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	return &copyConfig
}

func (this *logger) SetConfig(config *Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config = *config
	return nil
}

func (this *logger) UpdateConfig(fn func(*Config)) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := this.config
	fn(&copyConfig)
	this.config = copyConfig
	return nil
}

func (this *logger) Level() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level
}

func (this *logger) SetLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Level = level
}

func (this *logger) TimeLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.TimeLevel
}

func (this *logger) SetTimeLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.TimeLevel = level
}

func (this *logger) PanicLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.PanicLevel
}

func (this *logger) SetPanicLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.PanicLevel = level
}

func (this *logger) TrackLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.TrackLevel
}

func (this *logger) SetTrackLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.TrackLevel = level
}

func (this *logger) ExitLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.ExitLevel
}

func (this *logger) SetExitLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.ExitLevel = level
}

func (this *logger) Filter() Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Filter
}

func (this *logger) SetFilter(filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Filter = filter
}

func (this *logger) Prefix() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Prefix
}

func (this *logger) SetPrefix(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Prefix = ok
}

func (this *logger) Context() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Context
}

func (this *logger) SetContext(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Context = ok
}

func (this *logger) Mark() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Mark
}

func (this *logger) SetMark(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Mark = ok
}

func (this *logger) Limit() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Limit
}

func (this *logger) SetLimit(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Limit = ok
}
