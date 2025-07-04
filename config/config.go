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
	//1，2，3都是有关排行榜的key
	//1：增量哈希表的key
	RedisKeyMusicLikeIncrement string
	//2：排行榜zset的key（score->like_cnt,member->music_uuid)
	RedisKeyMusicSort string
	//3：存放具体数据哈希表的key ->哈希表->每一个元素是key,value结构 ，key为music_uuid，value为具体数据json)
	//因为考虑到增量的时候不能直接通过json快速查找到member，member是一个uuid的时候，才能通过增量哈希表维护的uuid快速找到
	//zset中的member，从而快速进行操作，那么其余数据只能维护在另外一个哈希结构中
	RedisKeyJsonData string
	RedisRankingsNum int64
}

var DefaultRedisKeyConfig = RedisKeyConfig{
	RedisKeyUserMusicLike: "like:user:%d:music:%s",
	//RedisKeyOfCollect:  "collect",
	RedisKeyMusicLikeCount: "like:count:music:%s",
	//维护排行版那个点赞增量集合的那个key（这个集合用于存放一个哈希，每一个哈希元素是音乐uuid与点赞增量数）
	RedisKeyMusicLikeIncrement: "like:increment:music",
	RedisKeyMusicSort:          "like:hot_file:zset",
	RedisKeyJsonData:           "like:hot_file:jsonhash:%s",
	RedisRankingsNum:           5,
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
