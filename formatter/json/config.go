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
type Config struct {
	// FileSegs specifies how many segments from last of the File field of a
	// Record will be formatted. The separator of segment is '/'.
	// If FileSegs is not specified, 0 is used which means all.
	FileSegs int
	// PkgSegs specifies how many segments from last of the Pkg field of a
	// Record will be formatted. The separator of segment is '/'.
	// If PkgSegs is not specified, 0 is used which means all.
	PkgSegs int
	// FuncSegs specifies how many segments from last of the Func field of a
	// Record will be formatted. The separator of segment is '.'.
	// If FuncSegs is not specified, 0 is used which means all.
	FuncSegs int
	// Omit specifies which fields of a Record will be omitted.
	Omit OmitBits
	// OmitEmpty specifies which fields of a Record will be omitted when they are
	// the zero value of their type.
	OmitEmpty OmitBits
	// MinBufSize is the initial size of the internal buf of a formatter.
	// MinBufSize must NOT be negative. If it is not specified, 384 is used.
	MinBufSize int
}

func (config *Config) setDefaults() {
	if config.MinBufSize == 0 {
		config.MinBufSize = 384
	}
}
