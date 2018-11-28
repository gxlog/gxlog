package unix

import "os"

const (
	DefaultPerm      = 0700
	DefaultOverwrite = true
)

type Config struct {
	Pathname  string
	Perm      os.FileMode
	Overwrite bool
}

func NewConfig(pathname string) *Config {
	return &Config{
		Pathname:  pathname,
		Perm:      DefaultPerm,
		Overwrite: DefaultOverwrite,
	}
}

func (this *Config) WithPerm(perm os.FileMode) *Config {
	this.Perm = perm
	return this
}

func (this *Config) WithOverwrite(ok bool) *Config {
	this.Overwrite = ok
	return this
}
