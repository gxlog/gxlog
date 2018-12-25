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

func (this *Config) WithFileSegs(segs int) *Config {
	this.FileSegs = segs
	return this
}

func (this *Config) WithPkgSegs(segs int) *Config {
	this.PkgSegs = segs
	return this
}

func (this *Config) WithFuncSegs(segs int) *Config {
	this.FuncSegs = segs
	return this
}

func (this *Config) WithOmit(bits OmitBits) *Config {
	this.Omit = bits
	return this
}

func (this *Config) WithOmitEmpty(bits OmitBits) *Config {
	this.OmitEmpty = bits
	return this
}

func (this *Config) WithMinBufSize(size int) *Config {
	this.MinBufSize = size
	return this
}
