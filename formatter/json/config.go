package json

type OmitBits int

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

type Config struct {
	FileSegs   int
	PkgSegs    int
	FuncSegs   int
	Omit       OmitBits
	OmitEmpty  OmitBits
	MinBufSize int
}

func NewConfig() *Config {
	return &Config{
		MinBufSize: 384,
	}
}

func (cfg *Config) WithFileSegs(segs int) *Config {
	cfg.FileSegs = segs
	return cfg
}

func (cfg *Config) WithPkgSegs(segs int) *Config {
	cfg.PkgSegs = segs
	return cfg
}

func (cfg *Config) WithFuncSegs(segs int) *Config {
	cfg.FuncSegs = segs
	return cfg
}

func (cfg *Config) WithOmit(bits OmitBits) *Config {
	cfg.Omit = bits
	return cfg
}

func (cfg *Config) WithOmitEmpty(bits OmitBits) *Config {
	cfg.OmitEmpty = bits
	return cfg
}

func (cfg *Config) WithMinBufSize(size int) *Config {
	cfg.MinBufSize = size
	return cfg
}
