package music

import (
	"MyChat/common/mysql"
	myredis "MyChat/common/redis"
	"MyChat/config"
	"MyChat/controller"
	"MyChat/model"
	"MyChat/utils/file"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var ctx = context.Background()

func IsExistMusicFile(user_id int64, file_id string) bool {
	pattern := fmt.Sprintf(config.GetConfig().MusicFilePath+"/%d/%s*", user_id, file_id)
	log.Println("file_path is ", pattern)
	return file.IsExistFile(pattern)
}

func UploadMusicFile(uuid, music_name, file_path string, user_id, file_size int64, isupload int64) (*model.MusicFile, bool) {
	if music_file, err := mysql.InsertMusicFile(&model.MusicFile{
		UUID:      uuid,
		UserID:    user_id,
		MusicName: music_name,
		FilePath:  file_path,
		IsUpload:  isupload,
		FileSize:  file_size,
	}); err != nil {
		log.Println(err.Error())
		return nil, false
	} else {
		return music_file, true
	}

}

// 从redis中获取前cnt个元素，如果没有对应的键值对，那么就直接从mysql中一次性加载所有数据到redis中
func GetTopInformation(cnt int64) ([]controller.MusicDetail, bool) {
	var res []controller.MusicDetail

	key := myredis.GenerateMusicLikeIncrementKey()
	//1:查找key值对应的zset是否存在   获取 ZSET 中前 cnt 个元素，并按照点赞数排序
	topMusicUUIDs, err := myredis.Rdb.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:   "-inf", // 无下限，获取所有数据
		Max:   "+inf", // 无上限
		Count: cnt,    // 获取前 cnt 个元素
	}).Result()
	//2：不存在：从mysql中加载音乐点赞数到redis中指定的zset中

	//3:从redis中读取对应的信息，并返回
	return res, true
}
