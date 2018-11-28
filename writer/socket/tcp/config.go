package tcp

import "errors"

type Config struct {
	Addr string
}

func NewConfig(addr string) *Config {
	return &Config{
		Addr: addr,
	}
}

func (this *Config) Check() error {
	if this.Addr == "" {
		return errors.New("Config.Addr should not be empty")
	}
	return nil
}
