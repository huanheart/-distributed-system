package user

import (
	"MyChat/common/mysql"
	"MyChat/common/rabbitmq"
	myredis "MyChat/common/redis"
	"MyChat/config"
	"MyChat/model"
	"MyChat/utils"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"
)

const (
	CodeMsg     = "MyChat验证码如下(验证码仅限于2分钟有效): "
	UserNameMsg = "MyChat的账号如下，请保留好，后续可以用账号/邮箱登录 "
)

var ctx = context.Background()

func SendCaptcha(email, code, msg string) error {
	m := gomail.NewMessage()

	// 发件人
	m.SetHeader("From", config.GetConfig().EmailConfig.Email)
	// 收件人
	m.SetHeader("To", email)
	// 主题
	m.SetHeader("Subject", "来自MyChat的信息")
	// 正文内容（纯文本形式，也可以用 text/html）
	m.SetBody("text/plain", msg+" "+code)

	// 配置 SMTP 服务器和授权码,587：是 SMTP 的明文/STARTTLS 端口号
	d := gomail.NewDialer("smtp.qq.com", 587, config.GetConfig().EmailConfig.Email, config.GetConfig().EmailConfig.Authcode)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("DialAndSend err %v:\n", err)
		return err
	}
	fmt.Printf("send mail success\n")
	return nil
}

func Register(username, email, password string) (*model.User, bool) {
	if user, err := mysql.InsertUser(&model.User{
		Email:    email,
		Name:     username,
		Username: username,
		Password: utils.MD5(password),
	}); err != nil {
		return nil, false
	} else {
		return user, true
	}
}

// 查询一个用户对这首歌的状态记录，先查询redis，如果redis中没有，那么就查询mysql
// 返回值0表示未点赞，1表示已点赞
func GetUserStatusOnFile(user_id int64, file_id string) (int64, bool) {
	//1：查看redis当前用户对这首音乐的点赞状态的key值
	key := myredis.GenerateLikeKey(user_id, file_id)

	//2:判断redis中是否有这个键值对(没有就从mysql中加载，否则直接返回redis中的数据）
	val, err := myredis.Rdb.Get(ctx, key).Result()
	if err == nil {
		if status, err := strconv.ParseInt(val, 10, 64); err == nil {
			return status, true
		}
		return 0, false
	} else if err == redis.Nil {
		//此时需要从mysql中加载对应状态到redis中
		reaction, err := mysql.GetMusicReaction(user_id, file_id)
		if err == nil && reaction == nil {
			//说明此时用户一定是对其进行点赞操作，即当前是未点赞状态（因为没有这行数据）
			return 0, true
		} else if err == nil {
			return reaction.Action, true
		}
		log.Println("mysql查询出错: " + err.Error())
		//否则就是mysql查询的时候出现问题了
		return 0, false
	}
	return 0, false
}

// 2：更新redis状态为LikeStatus相反的状态，并开启消息队列，将其mysql中对应reaction表中（插入或更新一条数据，更新现在是点赞还是取消点赞状态）
func ChangeOppositeState(user_id int64, file_id string, status int64) bool {
	key := myredis.GenerateLikeKey(user_id, file_id)
	//说明一开始是未点赞状态，现在需要变成点赞状态
	if status == 0 {
		myredis.Rdb.Set(ctx, key, 1, 0)
		//消息队列更新mysql
		//这里不关心Likecnt
		data := rabbitmq.GenerateLikeMQParam(user_id, 1, 0, file_id)
		rabbitmq.RMQUpdateAction.Publish(data)
		return true
	}
	//未点赞->点赞
	myredis.Rdb.Set(ctx, key, 0, 0)
	//消息队列更新mysql
	//需要将action表改成未点赞状态,这里不关心Likecnt
	data := rabbitmq.GenerateLikeMQParam(user_id, 0, 0, file_id)
	rabbitmq.RMQUpdateAction.Publish(data)
	return true
}

// 3：查看当前这首歌被点赞的总记录，先查询redis中是否含有，如果没有，mysql直接加载到redis中（此时mysql一定会有，因为歌曲是一定有的）
func GetFileLike(file_id string) (int64, bool) {
	var LikeCnt int64
	//查询redis是否有
	key := myredis.GenerateMusicCountKey(file_id)
	val, err := myredis.Rdb.Get(ctx, key).Result()
	if err == nil {
		LikeCnt, _ = strconv.ParseInt(val, 10, 64)
		return LikeCnt, true
	} else if err == redis.Nil {
		//需要从mysql中加载了,不出意外是不会出现music文件不存在的情况
		musicfile, err := mysql.GetMusicfileByFileId(file_id)
		if err != nil {
			log.Println("GetFileLike mysql解析错误原因: " + err.Error())
			return 0, false
		}
		LikeCnt = musicfile.LikeCount
		return LikeCnt, true
	}
	log.Println("GetFileLike redis解析错误原因: " + err.Error())
	//否则说明redis出错了
	return 0, false
}

// 4:更新redis当前点赞的总记录，并开启消息队列，将mysql中music_file表中的总点赞数进行+1/-1
func ChangeOppositeLikeCnt(status int64, LikeCnt int64, file_id string) bool {
	//根据LikeCnt做相反的操作，假设当前LikeStatus是未点赞状态，那么需要将redis中LikeCnt+1，反之
	key := myredis.GenerateMusicCountKey(file_id)
	if status == 0 {
		myredis.Rdb.Set(ctx, key, LikeCnt+1, 0)
		//放入消息队列中，更新mysql状态为LikeCnt+1
		data := rabbitmq.GenerateLikeMQParam(0, 0, LikeCnt+1, file_id)
		rabbitmq.RMQUpdateLikeCount.Publish(data)
		return true
	}
	myredis.Rdb.Set(ctx, key, LikeCnt-1, 0)
	//放入消息队列中，更新mysql状态为LikeCnt-1
	data := rabbitmq.GenerateLikeMQParam(0, 0, LikeCnt-1, file_id)
	rabbitmq.RMQUpdateLikeCount.Publish(data)
	return true
}
