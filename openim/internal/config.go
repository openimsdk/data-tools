package internal

import (
	"fmt"
	"github.com/openimsdk/data-tools/utils"
)

type Config struct {
	MySQL    utils.MySQLConfig `yaml:"mysql"`
	Mongo    utils.MongoConfig `yaml:"mongo"`
	S3Engine string            `yaml:"s3_engine"`
}

func (c Config) Check() error {
	if c.S3Engine == "" {
		return fmt.Errorf("config s3_engine not config")
	}
	return nil
}
