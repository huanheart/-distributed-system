package music

import (
	"MyChat/common/mysql"
	myredis "MyChat/common/redis"
	"MyChat/config"
	"MyChat/controller"
	"MyChat/model"
	"MyChat/utils"
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

// 如果没有对应的键值对，那么就直接从mysql中一次性加载所有数据到redis中(这里默认数据库是肯定存在数据的）
func LoadTopDataToRedis() bool {
	key := myredis.GenerateMusicLikeHotSortKey()
	//1:查找key值对应的zset是否存在
	zsetResults, _ := myredis.Rdb.ZRevRangeWithScores(ctx, key, 0, 0).Result()
	//2：不存在（切片为空）：从mysql中加载音乐点赞数到redis中指定的zset中
	if len(zsetResults) == 0 {
		//获取前cnt元素的信息
		musicFiles, err := mysql.GetTopAllFromMysql()
		if err != nil {
			log.Println("GetTopInformation mysql error: " + err.Error())
			return false
		}
		//遍历musicFiles并进行缓存操作
		var zMembers []*redis.Z
		for _, music := range musicFiles {
			//存放到zset中
			zMembers = append(zMembers, &redis.Z{
				Score:  float64(music.LikeCount),
				Member: music.UUID,
			})
			//存放到infohash中
			info_key := myredis.GenerateMusicJsonHashKey(music.UUID)
			_, err := myredis.Rdb.HMSet(ctx, info_key, map[string]interface{}{
				"file_path":  music.FilePath,
				"music_name": music.MusicName,
				"like_count": music.LikeCount,
			}).Result()
			if err != nil {
				log.Println("HMSet failed for", music.UUID, ":", err)
			}
		}
		_, err = myredis.Rdb.ZAdd(ctx, key, zMembers...).Result()
		if err != nil {
			log.Println("ZAdd failed: " + err.Error())
		}
	}
	return true
}

// 从redis中获取前cnt个元素
func GetTopInformation(cnt int64) ([]controller.MusicDetail, bool) {
	var res []controller.MusicDetail
	key := myredis.GenerateMusicLikeHotSortKey()
	//1. 获取 ZSet 中前 cnt 个 uuid（按点赞数排序）
	zsetResults, _ := myredis.Rdb.ZRevRangeWithScores(ctx, key, 0, cnt-1).Result()
	//2:遍历获取对应的哈希表，并取出数据放入到res中
	for _, z := range zsetResults {
		uuid := z.Member.(string)
		info_key := myredis.GenerateMusicJsonHashKey(uuid)

		info, err := myredis.Rdb.HGetAll(ctx, info_key).Result()
		if err != nil {
			log.Println("HGetAll error for uuid", uuid, ":", err)
			return nil, false
		}
		if len(info) == 0 {
			continue
		}
		var detail controller.MusicDetail
		detail.FilePath = utils.GetHttpPath(info["file_path"])
		detail.LikeCount = int64(z.Score) // 分数就是点赞数
		res = append(res, detail)
	}
	return res, true
}
