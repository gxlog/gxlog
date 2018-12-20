package tcp

type Config struct {
	Addr string
}

func NewConfig(addr string) *Config {
	return &Config{
		Addr: addr,
	}
}
