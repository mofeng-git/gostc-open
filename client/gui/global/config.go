package global

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

var Config config

type config struct {
	Address string `yaml:"address"`
}

func LoadConfig(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		Config = config{
			Address: "0.0.0.0:18080",
		}
		marshal, _ := yaml.Marshal(Config)
		if err := os.WriteFile(path, marshal, 0644); err != nil {
			return err
		}
		return nil
	}
	if stat.IsDir() {
		return errors.New(path + " is not file")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, &Config)
}

func SaveConfig(path string) error {
	marshal, _ := yaml.Marshal(Config)
	return os.WriteFile(path, marshal, 0644)
}
