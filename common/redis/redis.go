package redis

import (
	"MyChat/config"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var Rdb *redis.Client

func Init() {
	conf := config.GetConfig()
	host := conf.RedisConfig.RedisHost
	port := conf.RedisConfig.RedisPort
	password := conf.RedisConfig.RedisPassword
	db := conf.RedisDb
	addr := host + ":" + strconv.Itoa(port)

	Rdb = redis.NewClient(&redis.Options{
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
	return Rdb.Set(ctx, key, captcha, expire).Err()
}

func CheckCaptchaForEmail(ctx context.Context, email, userInput string) (bool, error) {
	key := redisCaptcha(email)

	// 从 Redis 获取验证码
	storedCaptcha, err := Rdb.Get(ctx, key).Result()
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

//func SetLike(ctx context.Context, key string, val string) bool {
//	if val == "0" {
//		//说明未点赞，将其转化成点赞，并查询是否含有music字符串，并进行更新数量更新
//		Rdb.Set(ctx, key, "1", 0)
//	} else {
//		Rdb.Set(ctx, key, "0", 0)
//	}
//	return true
//}
//
//// 对音乐通过val进行+1 -1操作
//func UpdateLikeByVal(ctx context.Context, key string, val string) bool {
//	//说明要将其+1
//	if val == "0" {
//		return UpdateLike(ctx, key, 1)
//	} else {
//		return UpdateLike(ctx, key, -1)
//	}
//}
//func UpdateLike(ctx context.Context, key string, delta int64) bool {
//	err := Rdb.IncrBy(ctx, key, int64(delta)).Err()
//	return err == nil
//}
