package tcp

// A Config is used to configure a tcp socket writer.
// A Config should be created with NewConfig.
type Config struct {
	Addr string
}

// NewConfig creates a new Config.
func NewConfig(addr string) *Config {
	return &Config{
		Addr: addr,
	}
}
