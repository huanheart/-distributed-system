package redis

import (
	"MyChat/config"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var rdb *redis.Client

func Init() {
	conf := config.GetConfig()
	host := conf.RedisConfig.RedisHost
	port := conf.RedisConfig.RedisPort
	password := conf.RedisConfig.RedisPassword
	db := conf.RedisDb
	addr := host + ":" + strconv.Itoa(port)

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

//github.com/go-redis/redis/v8
//key:特定邮箱-> 验证码

func redisCaptcha(str string) string {
	return "captcha:" + str
}

func SetCaptchaForEmail(ctx context.Context, email, captcha string) error {
	key := redisCaptcha(email)
	expire := 2 * time.Minute
	return rdb.Set(ctx, key, captcha, expire).Err()
}

func CheckCaptchaForEmail(ctx context.Context, email, userInput string) (bool, error) {
	key := redisCaptcha(email)

	// 从 Redis 获取验证码
	storedCaptcha, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Redis 中没有这个 key，说明验证码不存在或已过期
			return false, nil
		}
		return false, err // 其他 Redis 错误
	}

	// 比较验证码是否一致
	return storedCaptcha == userInput, nil
}
