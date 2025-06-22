package user

import (
	"MyChat/common/mysql"
	"MyChat/common/redis"
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
	if ok, _ := redis.CheckCaptchaForEmail(ctx, email, captcha); !ok {

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

	if err := redis.SetCaptchaForEmail(ctx, email, code); err != nil {
		return false
	}

	//2:再进行远程发送
	if err := user.SendCaptcha(email, code, user.CodeMsg); err != nil {
		return false
	}

	return true
}
