package user

import (
	"MyChat/common/mysql"
	myredis "MyChat/common/redis"
	"MyChat/dao/user"
	"MyChat/model"
	"MyChat/utils"
	"context"
	"gorm.io/gorm"
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

// 最多产生多少次redis和mysql的当前消耗？
// 4次redis + 2次mysql
// MQ减少了多少次mysql操作？
// 2次以上
func HandleLike(user_id int64, file_id string) (int64, int64, bool) {
	var LikeCnt int64
	var LikeStatus int64
	var ok bool
	// 1：查询一个用户对这首歌的状态记录，先查询redis，如果redis中没有，那么就加载mysql数据到redis中，如果mysql中也没有
	// 说明此时一定是未点赞情况,那么加载对应未点赞状态到redis中
	if LikeStatus, ok = user.GetUserStatusOnFile(user_id, file_id); !ok {
		return 0, 0, false
	}
	// 2：更新redis状态为LikeStatus相反的状态，并开启消息队列，将其mysql中对应reaction表中（插入或更新一条数据，更新现在是点赞还是取消点赞状态）
	user.ChangeOppositeState(user_id, file_id, LikeStatus)

	// 3：查看当前这首歌被点赞的总记录，先查询redis中是否含有，如果没有，mysql直接加载到redis中（此时mysql一定会有，因为歌曲是一定有的）
	LikeCnt, _ = user.GetFileLike(file_id)

	//4:更新redis当前点赞的总记录，并开启消息队列，将mysql中music_file表中的总点赞数进行+1/-1
	user.ChangeOppositeLikeCnt(LikeStatus, LikeCnt, file_id)

	// 5：最后还要将当前状态进行反转,返回当前总点赞数以及当前用户对这个音乐的点赞状态
	if LikeStatus == 0 {
		LikeStatus = 1
		LikeCnt = LikeCnt + 1
	} else {
		LikeStatus = 0
		LikeCnt = LikeCnt - 1
	}
	return LikeCnt, LikeStatus, true
}
