package tcp

// A Config is used to configure a tcp socket writer.
type Config struct {
	// If Tag is not specified, "localhost:9999" is used.
	Addr string
}

func (config *Config) setDefaults() {
	if config.Addr == "" {
		config.Addr = "localhost:9999"
	}
}
