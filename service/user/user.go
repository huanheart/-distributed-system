package user

import (
	"MyChat/common/mysql"
	myredis "MyChat/common/redis"
	"MyChat/dao/user"
	"MyChat/model"
	"MyChat/utils"
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
)

var ctx = context.Background()

// 拿邮箱或者账号来判断都可以
func IsExistUser(account string) (bool, *model.User) {
	user, err := mysql.GetUserByEmail(account)
	if err == nil {
		return true, user
	}
	user, err = mysql.GetUserByUsername(account)

	if err == gorm.ErrRecordNotFound || user == nil {
		return false, nil
	}

	return true, user
}

// 拿邮箱注册一个账户
func Register(email, password, captcha string) (bool, *model.User) {
	var ok bool
	var user_ *model.User
	//2:开始注册
	//2.0:从redis中验证验证码是否有效
	if ok, _ := myredis.CheckCaptchaForEmail(ctx, email, captcha); !ok {

		return false, nil
	}
	//2.1：生成11位的账号
	username := utils.GetRandomNumbers(11)
	//2.2：注册到数据库中

	if user_, ok = user.Register(username, email, password); !ok {
		return false, nil
	}
	//2.3：将账号一并发送到对应邮箱上去，后续可以选择账号和邮箱登录
	if err := user.SendCaptcha(email, username, user.UserNameMsg); err != nil {
		return false, nil
	}

	return true, user_
}

// 往指定邮箱发送验证码
// 分为以下任务：
// 1：先存放redis
// 2：再进行远程发送
func SendCaptcha(email string) bool {
	code := utils.GetRandomNumbers(6)
	//1:先存放redis

	if err := myredis.SetCaptchaForEmail(ctx, email, code); err != nil {
		return false
	}

	//2:再进行远程发送
	if err := user.SendCaptcha(email, code, user.CodeMsg); err != nil {
		return false
	}

	return true
}

func HandleLike(user_id int64, file_id string) (int64, int64, bool) {
	var LikeCnt int64
	var LikeStatus int64

	//1：查看redis当前用户对这首音乐的点赞状态,查询是否存在redis中
	key := myredis.GenerateLikeKey(user_id, file_id)
	//2:判断redis中是否有这个键值对
	val, err := myredis.Rdb.Get(ctx, key).Result()
	//说明有这个键值对，获取它是点赞还是未点赞
	if err == nil {
		//将点赞->非点赞 / 非点赞->点赞
		if ok := myredis.SetLike(ctx, key, val); !ok {
			return LikeCnt, LikeStatus, false
		}
		//获取key值，用于查询当前音乐播放量的总数
		key = myredis.GenerateMusicCountKey(file_id)
		_, err := myredis.Rdb.Get(ctx, key).Result()
		if err == nil {
			//将对应的值加1或者-1
			myredis.UpdateLikeByVal(ctx, key, val)
		} else if err == redis.Nil {
			//查询mysql的Music_file表，获取该音乐被点赞总数并放入到redis中
			//  musicfile极大概率都是查找的到的，这里不考虑异常意外的情况
			musicfile, _ := mysql.GetMusicfileByFileId(file_id)
			LikeCnt = musicfile.LikeCount
			//存放到redis中
			myredis.UpdateLike(ctx, key, LikeCnt+1)
			// todo:放入到消息队列中，需要更新mysql的字段,这里异步操作（两张表都需要）
			// todo 2: 将music_file表中的LikeCount进行+1操作

		} else {
			return LikeCnt, LikeStatus, false
		}
		//todo: 1:将music_reaction表中的Action更改为val相反的值

	} else if err == redis.Nil {
		//说明没有这个键值对，需要从mysql中查询MusicReaction表获取当前用户对这个歌曲的点赞状态，并存入redis
		reaction, err := mysql.GetMusicReaction(user_id, file_id)
		if reaction == nil && err == nil {
			//1：直接将其存放到redis中，此时一定是点赞操作,固然从0->1
			if ok := myredis.SetLike(ctx, key, "0"); !ok {
				return LikeCnt, LikeStatus, false
			}
			//todo:这里往Reaction表中插入一条数据(即user_id对file_id进行1（点赞操作）

			//然后开始查询数据
			//获取key值，用于查询当前音乐播放量的总数
			key = myredis.GenerateMusicCountKey(file_id)
			_, err := myredis.Rdb.Get(ctx, key).Result()
			if err == nil {
				//将对应的值加1或者-1
				myredis.UpdateLikeByVal(ctx, key, val)
			} else if err == redis.Nil {
				//查询mysql的Music_file表，获取该音乐被点赞总数并放入到redis中
				//  musicfile极大概率都是查找的到的，这里不考虑异常意外的情况
				musicfile, _ := mysql.GetMusicfileByFileId(file_id)
				LikeCnt = musicfile.LikeCount
				//存放到redis中
				myredis.UpdateLike(ctx, key, LikeCnt+1)
				// todo:放入到消息队列中，需要更新mysql的LikeCount字段,这里异步操作
			} else {
				return LikeCnt, LikeStatus, false
			}

		} else if reaction != nil {
			//直接将其值缓存到redis中,1->0 0->1
			myredis.SetLike(ctx, key, strconv.Itoa(int(reaction.Action)))
			//todo:更新reaction表中reaction.Action字段，变成相反的值即可
			//然后开始查询数量
			key = myredis.GenerateMusicCountKey(file_id)
			_, err := myredis.Rdb.Get(ctx, key).Result()
			if err == nil {
				//将对应的值加1或者-1
				myredis.UpdateLikeByVal(ctx, key, val)
			} else if err == redis.Nil {
				//查询mysql的Music_file表，获取该音乐被点赞总数并放入到redis中
				//  musicfile极大概率都是查找的到的，这里不考虑异常意外的情况
				musicfile, _ := mysql.GetMusicfileByFileId(file_id)
				LikeCnt = musicfile.LikeCount
				//存放到redis中
				myredis.UpdateLike(ctx, key, LikeCnt+1)
				// todo:放入到消息队列中，需要更新mysql的LikeCount字段,这里异步操作
			} else {
				return LikeCnt, LikeStatus, false
			}

		} else {
			log.Println(err.Error())
			return LikeCnt, LikeStatus, false
		}

	} else {
		log.Println("redis Get操作失败 " + err.Error())
		return LikeCnt, LikeStatus, false
	}

	return LikeCnt, LikeStatus, true
}
