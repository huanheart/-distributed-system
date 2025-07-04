package redis

import (
	"MyChat/config"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

var Rdb *redis.Client

var ctx = context.Background()

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

	//定时协程的封装,调用UpdateRedisCache
	//用于定时更新redis中的缓存的函数(将key中所有的增量全部更新到zset与info_hash中）
	go func() {
		ticker := time.NewTicker(1 * time.Minute) // 每小时
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("[定时任务] 开始执行 UpdateRedisCache...")
				UpdateRedisCache()
			}
		}
	}()

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

// 用于定时更新redis中的缓存的函数(将key中所有的增量全部更新到zset与info_hash中）
//
// 用增量模式
func UpdateRedisCache() {
	hashKey := GenerateMusicLikeIncrementKey()
	zsetKey := GenerateMusicLikeHotSortKey()
	// 1. 获取整个哈希表的所有字段和值（uuid -> 增量）
	increments, err := Rdb.HGetAll(ctx, hashKey).Result()
	if err != nil || len(increments) == 0 {
		return
	}

	for uuid, incrStr := range increments {
		incr, _ := strconv.ParseInt(incrStr, 10, 64)
		if incr == 0 {
			continue // 不更新
		}
		//更新 ZSet 中的分数
		_, err = Rdb.ZIncrBy(ctx, zsetKey, float64(incr), uuid).Result()
		if err != nil {
			log.Printf("UpdateRedisCache - ZIncrBy failed for %s: %v\n", uuid, err)
			continue
		}
		_, err = Rdb.HDel(ctx, hashKey, uuid).Result()
		if err != nil {
			log.Printf("UpdateRedisCache - HDel failed for %s: %v\n", uuid, err)
		}
		//内部更新info_hash中的like_count，更改它的增量
		infohashKey := GenerateMusicJsonHashKey(uuid)
		_, err = Rdb.HIncrBy(ctx, infohashKey, "like_count", incr).Result()
		if err != nil {
			log.Printf("UpdateRedisCache - HIncrBy infohash failed for %s: %v\n", uuid, err)
		}

	}

}

func UpdateLikeIncrement(musicUUID string, increment int64) bool {
	key := GenerateMusicLikeIncrementKey() // 你定义的哈希 key，比如 "like:increment:music"
	//如果map[musicUUID]不存在，那么会自动创建对应musicUUID为内部key，value为0的关系
	_, err := Rdb.HIncrBy(ctx, key, musicUUID, increment).Result()
	if err != nil {
		log.Println("UpdateLikeIncrement failed:", err)
		return false
	}
	return true
}

func AddOneLikeIncrement(musicUUID string) bool {
	return UpdateLikeIncrement(musicUUID, 1)
}

func SubOneLikeIncrement(musicUUID string) bool {
	return UpdateLikeIncrement(musicUUID, -1)
}
