package socket

type Config struct {
	Network string
	Bind    string
}

func NewConfig(network, bind string) *Config {
	return &Config{
		Network: network,
		Bind:    bind,
	}
}
