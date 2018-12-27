package unix

import (
	"os"
)

type Config struct {
	Pathname  string
	Perm      os.FileMode
	Overwrite bool
}

func NewConfig(pathname string) *Config {
	return &Config{
		Pathname:  pathname,
		Perm:      0700,
		Overwrite: true,
	}
}

func (cfg *Config) WithPerm(perm os.FileMode) *Config {
	cfg.Perm = perm
	return cfg
}

func (cfg *Config) WithOverwrite(ok bool) *Config {
	cfg.Overwrite = ok
	return cfg
}
