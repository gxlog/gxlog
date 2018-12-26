package gxlog

// The Flag defines the flag type of Logger.
type Flag int

// All available flags of Logger here. If a flag is set, its corresponding
// feature is enabled.
const (
	Prefix Flag = 0x01 << iota
	Contexts
	DynamicContexts // to be effective, the flag Contexts must also be set
	Mark
	Limit // corresponding to limiting by count and limiting by time
)

// The Filter type defines a function type which is as the log filter.
// Do NOT call methods of the Logger within a filter, or it will deadlock.
type Filter func(*Record) bool

// A Config is used to configure a Logger.
// A Config should be created with NewConfig.
type Config struct {
	// Level is the level of Logger, logs with lower level will be omitted.
	Level Level
	// TimeLevel is the level of a log when Time or Timef is called.
	TimeLevel Level
	// PanicLevel is the level of a log when Panic or Panicf is called.
	PanicLevel Level
	// TrackLevel is the auto backtracking level of Logger.
	// If the level of a log is NOT lower than the TrackLevel, the stack of
	// the current goroutine will also be output.
	TrackLevel Level
	// ExitLevel is the auto exiting level of Logger.
	// If the level of a log is NOT lower than the ExitLevel, the Logger
	// will call os.Exit after outputs log.
	ExitLevel Level
	// Filter is the log filter of Logger. If it is not nil and it returns
	// false, the log will be omitted.
	Filter Filter
	Flags  Flag
}

// NewConfig creates a new Config. By default, the Level is Trace, the
// TimeLevel is Trace, the PanicLevel is Fatal, the TrackLevel is Fatal,
// the ExitLevel is Off, the Filter is nil and all flags are set.
func NewConfig() *Config {
	return &Config{
		Level:      Trace,
		TimeLevel:  Trace,
		PanicLevel: Fatal,
		TrackLevel: Fatal,
		ExitLevel:  Off,
		Flags:      Prefix | Contexts | DynamicContexts | Mark | Limit,
	}
}

// WithLevel sets the Level of the Config and returns it.
func (this *Config) WithLevel(level Level) *Config {
	this.Level = level
	return this
}

// WithTimeLevel sets the TimeLevel of the Config and returns it.
func (this *Config) WithTimeLevel(level Level) *Config {
	this.TimeLevel = level
	return this
}

// WithPanicLevel sets the PanicLevel of the Config and returns it.
func (this *Config) WithPanicLevel(level Level) *Config {
	this.PanicLevel = level
	return this
}

// WithTrackLevel sets the TrackLevel of the Config and returns it.
func (this *Config) WithTrackLevel(level Level) *Config {
	this.TrackLevel = level
	return this
}

// WithExitLevel sets the ExitLevel of the Config and returns it.
func (this *Config) WithExitLevel(level Level) *Config {
	this.ExitLevel = level
	return this
}

// WithFilter sets the Filter of the Config and returns it.
func (this *Config) WithFilter(filter Filter) *Config {
	this.Filter = filter
	return this
}

// WithFlags sets the Flags of the Config and returns it.
func (this *Config) WithFlags(flags Flag) *Config {
	this.Flags = flags
	return this
}

// WithEnabled enables flags of the Config and returns it.
func (this *Config) WithEnabled(flags Flag) *Config {
	this.Flags |= flags
	return this
}

// WithDisabled disables flags of the Config and returns it.
func (this *Config) WithDisabled(flags Flag) *Config {
	this.Flags &^= flags
	return this
}
