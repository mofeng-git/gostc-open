package configs

type Sqlite struct {
	File     string `yaml:"file"`
	LogLevel string `yaml:"log-level"`
}
