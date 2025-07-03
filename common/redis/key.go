package redis

import (
	"MyChat/config"
	"fmt"
)

// 用于判断某个用户是否对某个音乐进行点赞
func GenerateLikeKey(user_id int64, file_id string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.RedisKeyUserMusicLike, user_id, file_id)
}

// 统计某个音乐的总点赞次数
func GenerateMusicCountKey(file_id string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.RedisKeyMusicLikeCount, file_id)
}

func GenerateMusicLikeIncrementKey() string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.RedisKeyMusicLikeIncrement)
}
