package unix

import (
	"os"
	"strconv"
)

// A Config is used to configure a unix domain socket writer.
type Config struct {
	// Pathname is the pathname of the socket file that will be created.
	// Shell expansion is NOT supported.
	// If Pathname is not specified, "/tmp/gxlog/<pid>" is used.
	Pathname string
	// Perm is permission of the socket file that will be created.
	// If Perm is not specified, 0700 is used.
	Perm os.FileMode
	// NoOverwrite specifies NOT to overwrite a existing socket file.
	NoOverwrite bool
}

func (config *Config) setDefaults() {
	if config.Pathname == "" {
		config.Pathname = "/tmp/gxlog/" + strconv.Itoa(os.Getpid())
	}
	if config.Perm == 0 {
		config.Perm = 0700
	}
}
