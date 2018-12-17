package unix

import (
	"errors"
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

func (this *Config) WithPerm(perm os.FileMode) *Config {
	this.Perm = perm
	return this
}

func (this *Config) WithOverwrite(ok bool) *Config {
	this.Overwrite = ok
	return this
}

func (this *Config) Check() error {
	if this.Pathname == "" {
		return errors.New("Config.Pathname should not be empty")
	}
	return nil
}
