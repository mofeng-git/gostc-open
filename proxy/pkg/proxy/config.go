package proxy

type Config struct {
	HTTPAddr  string                  `yaml:"http-addr"`
	HTTPSAddr string                  `yaml:"https-addr"`
	Default   DomainConfig            `yaml:"default"`
	Domains   map[string]DomainConfig `yaml:"domains"`
}

type DomainConfig struct {
	Target string `yaml:"target"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
}
