package text

const DefaultHeader = "{{time}} {{level}} {{pathname}}:{{line}} {{func}} {{msg}}\n"

type Config struct {
	Header string
}

func NewConfig() *Config {
	return &Config{
		Header: DefaultHeader,
	}
}

func (this *Config) WithHeader(header string) *Config {
	this.Header = header
	return this
}
