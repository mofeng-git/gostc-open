package configs

import (
	"time"
)

type Config struct {
	Address   string        `yaml:"address"`
	AuthKey   string        `yaml:"auth-key"`
	AuthExp   time.Duration `yaml:"auth-exp"`
	AuthRenew time.Duration `yaml:"auth-renew"`
	DbType    string        `yaml:"db-type"`
	Sqlite    Sqlite        `yaml:"sqlite"`
	Mysql     Mysql         `yaml:"mysql"`
}
