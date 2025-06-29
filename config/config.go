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
	RedisPort     int    `toml:"port"`
	RedisDb       int    `toml:"db"`
	RedisHost     string `toml:"host"`
	RedisPassword string `toml:"password"`
}

type MysqlConfig struct {
	MysqlPort         int    `toml:"port"`
	MysqlHost         string `toml:"host"`
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
	Port          int    `toml:"port"`
	AppName       string `toml:"appName"`
	Host          string `toml:"host"`
	MusicFilePath string `toml:"musicFilePath"`
	HttpFilePath  string `toml:"httpFilePath"`
	MusicFileIp   string `toml:"musicFileIp"`
}

type Rabbitmq struct {
	RabbitmqPort     int    `toml:"port"`
	RabbitmqHost     string `toml:"host"`
	RabbitmqUsername string `toml:"username"`
	RabbitmqPassword string `toml:"password"`
	RabbitmqVhost    string `toml:"vhost"`
}

type Config struct {
	EmailConfig `toml:"emailConfig"`
	RedisConfig `toml:"redisConfig"`
	MysqlConfig `toml:"mysqlConfig"`
	JwtConfig   `toml:"jwtConfig"`
	MainConfig  `toml:"mainConfig"`
	Rabbitmq    `toml:"rabbitmqConfig"`
}

type RedisKeyConfig struct {
	RedisKeyUserMusicLike string
	//RedisKeyOfCollect  string
	RedisKeyMusicLikeCount string
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	RedisKeyUserMusicLike: "like:user:%d:music:%s",
	//RedisKeyOfCollect:  "collect",
	RedisKeyMusicLikeCount: "like:count:music:%s",
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
