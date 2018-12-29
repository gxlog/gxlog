package unix

import (
	"os"
)

// A Config is used to configure a unix domain socket writer.
// A Config should be created with NewConfig.
type Config struct {
	// Pathname is the pathname of the socket file that will be created.
	// Shell expansion is NOT supported.
	Pathname string
	// Perm is permission of the socket file that will be created.
	Perm os.FileMode
	// Overwrite specifies whether to overwrite a existing socket file.
	Overwrite bool
}

// NewConfig creates a new Config. The default Perm is 0700 and the default
// Overwrite is true.
func NewConfig(pathname string) *Config {
	return &Config{
		Pathname:  pathname,
		Perm:      0700,
		Overwrite: true,
	}
}

// WithPerm sets the Perm of the Config and returns the Config.
func (cfg *Config) WithPerm(perm os.FileMode) *Config {
	cfg.Perm = perm
	return cfg
}

// WithOverwrite sets the Overwrite of the Config and returns the Config.
func (cfg *Config) WithOverwrite(ok bool) *Config {
	cfg.Overwrite = ok
	return cfg
}
