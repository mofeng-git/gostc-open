package configs

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Prefix   string `yaml:"prefix"`
	Extend   string `yaml:"extend"`
	LogLevel string `yaml:"log-level"`
}
