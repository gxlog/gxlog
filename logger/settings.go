package logger

import (
	"github.com/gxlog/gxlog/iface"
)

// Config returns the Config of the Logger.
func (log *Logger) Config() Config {
	log.lock.Lock()
	defer log.lock.Unlock()

	return *log.config
}

// SetConfig sets the config to the Logger.
func (log *Logger) SetConfig(config Config) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config = &config
}

// UpdateConfig calls the fn with the Config of the Logger, and then sets the
// returned Config to the Logger. The fn must NOT be nil.
//
// Do NOT call any method of the Logger within the fn, or it may deadlock.
func (log *Logger) UpdateConfig(fn func(Config) Config) {
	log.lock.Lock()
	defer log.lock.Unlock()

	config := fn(*log.config)
	log.config = &config
}

// Level returns the level of the Logger.
func (log *Logger) Level() iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level
}

// SetLevel sets the level of the Logger.
func (log *Logger) SetLevel(level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.Level = level
}

// TrackLevel returns the track level of the Logger.
func (log *Logger) TrackLevel() iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.TrackLevel
}

// SetTrackLevel sets the track level of the Logger.
func (log *Logger) SetTrackLevel(level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.TrackLevel = level
}

// ExitLevel returns the exit level of the Logger.
func (log *Logger) ExitLevel() iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.ExitLevel
}

// SetExitLevel sets the exit level of the Logger.
func (log *Logger) SetExitLevel(level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.ExitLevel = level
}

// TimeLevel returns the time level of the Logger.
func (log *Logger) TimeLevel() iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.TimeLevel
}

// SetTimeLevel sets the time level of the Logger.
func (log *Logger) SetTimeLevel(level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.TimeLevel = level
}

// PanicLevel returns the panic level of the Logger.
func (log *Logger) PanicLevel() iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.PanicLevel
}

// SetPanicLevel sets the panic level of the Logger.
func (log *Logger) SetPanicLevel(level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.PanicLevel = level
}

// Filter returns the filter of the Logger.
func (log *Logger) Filter() Filter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Filter
}

// SetFilter sets the filter of the Logger.
func (log *Logger) SetFilter(filter Filter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.Filter = filter
}

// Disabled returns the disabled flags of the Logger.
func (log *Logger) Disabled() Flag {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Disabled
}

// SetDisabled sets the disabled flags of the Logger.
func (log *Logger) SetDisabled(flags Flag) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.Disabled = flags
}

// Enable enables the flags of the Logger.
func (log *Logger) Enable(flags Flag) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.Disabled &^= flags
}

// Disable disables the flags of the Logger.
func (log *Logger) Disable(flags Flag) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.Disabled |= flags
}
