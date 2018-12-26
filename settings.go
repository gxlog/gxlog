package gxlog

// Config returns a copy of Config of the Logger.
func (this *Logger) Config() *Config {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := *this.config
	return &copyConfig
}

// SetConfig sets the copy of config to the Logger. The config must NOT be nil.
func (this *Logger) SetConfig(config *Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	copyConfig := *config
	this.config = &copyConfig
	return nil
}

// UpdateConfig will call fn with copy of the config of the Logger, and then
// sets copy of the returned config to the Logger. The fn must NOT be nil.
// Do NOT call methods of the Logger within fn, or it will deadlock.
func (this *Logger) UpdateConfig(fn func(Config) Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	config := fn(*this.config)
	this.config = &config
	return nil
}

// Level returns the level of the Logger.
func (this *Logger) Level() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level
}

// SetLevel sets the level of the Logger.
func (this *Logger) SetLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Level = level
}

// TimeLevel returns the time level of the Logger.
func (this *Logger) TimeLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.TimeLevel
}

// SetTimeLevel sets the time level of the Logger.
func (this *Logger) SetTimeLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.TimeLevel = level
}

// PanicLevel returns the panic level of the Logger.
func (this *Logger) PanicLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.PanicLevel
}

// SetPanicLevel sets the panic level of the Logger.
func (this *Logger) SetPanicLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.PanicLevel = level
}

// TrackLevel returns the track level of the Logger.
func (this *Logger) TrackLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.TrackLevel
}

// SetTrackLevel sets the track level of the Logger.
func (this *Logger) SetTrackLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.TrackLevel = level
}

// ExitLevel returns the exit level of the Logger.
func (this *Logger) ExitLevel() Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.ExitLevel
}

// SetExitLevel sets the exit level of the Logger.
func (this *Logger) SetExitLevel(level Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.ExitLevel = level
}

// Filter returns the filter of the Logger.
func (this *Logger) Filter() Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Filter
}

// SetFilter sets the filter of the Logger.
func (this *Logger) SetFilter(filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Filter = filter
}

// Flags returns the flags of the Logger.
func (this *Logger) Flags() Flag {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Flags
}

// SetFlags sets the flags of the Logger.
func (this *Logger) SetFlags(flags Flag) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Flags = flags
}

// Enable enables the flags of the Logger.
func (this *Logger) Enable(flags Flag) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Flags |= flags
}

// Disable disables the flags of the Logger.
func (this *Logger) Disable(flags Flag) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Flags &^= flags
}
