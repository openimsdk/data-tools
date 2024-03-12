package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLConfig struct {
	Address     []string `yaml:"address"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	Database    string   `yaml:"database"`
	MaxOpenConn int      `yaml:"maxOpenConn"`
	MaxIdleConn int      `yaml:"maxIdleConn"`
	MaxLifeTime int      `yaml:"maxLifeTime"`
}

func (c MySQLConfig) buildURI() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Address[0], c.Database)
}

func NewMySQL(conf MySQLConfig) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(conf.buildURI()), &gorm.Config{Logger: logger.Discard})
}
