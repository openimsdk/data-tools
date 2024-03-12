package internal

import (
	"github.com/openimsdk/data-tools/utils"
)

type Config struct {
	MySQL utils.MySQLConfig `yaml:"mysql"`
	Mongo utils.MongoConfig `yaml:"mongo"`
}

func (c Config) Check() error {
	return nil
}
