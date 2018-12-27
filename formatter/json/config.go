package json

// The OmitBits defines the flag type that is used to omit fields of a Record.
type OmitBits int

// All available flags here. If a flag is set, the corresponding field of a
// Record will be omitted.
const (
	Time OmitBits = 0x1 << iota
	Level
	File
	Line
	Pkg
	Func
	Msg
	Prefix
	Context
	Mark
	Aux = Prefix | Context | Mark
)

// A Config is used to configure a json formatter.
// A Config should be created with NewConfig.
type Config struct {
	// FileSegs specifies how many segments from last of the File field will be
	// formatted. The separator is '/'.
	FileSegs int
	// PkgSegs specifies how many segments from last of the Pkg field will be
	// formatted. The separator is '/'.
	PkgSegs int
	// FuncSegs specifies how many segments from last of the Func field will be
	// formatted. The separator is '.'.
	FuncSegs int
	// Omit specifies which fields will be omitted.
	Omit OmitBits
	// OmitEmpty specifies which fields will be omitted when they are the zero
	// value of their type.
	OmitEmpty OmitBits
	// MinBufSize is the initial size of the internal buf of a formatter.
	// MinBufSize must not be negative.
	MinBufSize int
}

// NewConfig creates a new Config. By default, the FileSegs, PkgSegs and FuncSegs
// are 0, the Omit and OmitEmpty are all unset, the MinBufSize is 384.
func NewConfig() *Config {
	return &Config{
		MinBufSize: 384,
	}
}

// WithFileSegs sets the FileSegs of the Config and returns it.
func (cfg *Config) WithFileSegs(segs int) *Config {
	cfg.FileSegs = segs
	return cfg
}

// WithPkgSegs sets the PkgSegs of the Config and returns it.
func (cfg *Config) WithPkgSegs(segs int) *Config {
	cfg.PkgSegs = segs
	return cfg
}

// WithFuncSegs sets the FuncSegs of the Config and returns it.
func (cfg *Config) WithFuncSegs(segs int) *Config {
	cfg.FuncSegs = segs
	return cfg
}

// WithOmit sets the Omit of the Config and returns it.
func (cfg *Config) WithOmit(bits OmitBits) *Config {
	cfg.Omit = bits
	return cfg
}

// WithOmitEmpty sets the OmitEmpty of the Config and returns it.
func (cfg *Config) WithOmitEmpty(bits OmitBits) *Config {
	cfg.OmitEmpty = bits
	return cfg
}

// WithMinBufSize sets the MinBufSize of the Config and returns it.
func (cfg *Config) WithMinBufSize(size int) *Config {
	cfg.MinBufSize = size
	return cfg
}
