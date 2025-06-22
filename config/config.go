package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type EmailConfig struct {
	Authcode string `toml:"authcode"`
	Email    string `toml:"email" `
}

type RedisConfig struct {
	RedisHost     string `toml:"host"`
	RedisPort     int    `toml:"port"`
	RedisPassword string `toml:"password"`
	RedisDb       int    `toml:"db"`
}

type MysqlConfig struct {
	MysqlHost         string `toml:"host"`
	MysqlPort         int    `toml:"port"`
	MysqlUser         string `toml:"user"`
	MysqlPassword     string `toml:"password"`
	MysqlDatabaseName string `toml:"databaseName"`
	MysqlCharset      string `toml:"charset"`
}

type JwtConfig struct {
	ExpireDuration int    `toml:"expire_duration"`
	Issuer         string `toml:"issuer"`
	Subject        string `toml:"subject"`
	Key            string `toml:"key"`
}

type MainConfig struct {
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
}

type Config struct {
	EmailConfig `toml:"emailConfig"`
	RedisConfig `toml:"redisConfig"`
	MysqlConfig `toml:"mysqlConfig"`
	JwtConfig   `toml:"jwtConfig"`
	MainConfig  `toml:"mainConfig"`
}

var config *Config

// InitConfig 初始化项目配置
func InitConfig() error {
	// 设置配置文件路径（相对于 main.go 所在的目录）
	if _, err := toml.DecodeFile("config/config.toml", config); err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		_ = InitConfig()
	}
	return config
}
