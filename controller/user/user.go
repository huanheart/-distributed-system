package user

import (
	"MyChat/common/code"
	"MyChat/controller"
	"MyChat/model"
	"MyChat/service/user"
	"MyChat/utils"
	"MyChat/utils/myjwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	//这里的Username可以是账号也可以是邮箱
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:password`
	}
	// omitempty当字段为空的时候，不返回这个东西
	LoginResponse struct {
		controller.Response
		UserID int64  `json:"user_id,omitempty"`
		Token  string `json:"token,omitempty"`
	}
	//验证码由后端生成，存放到redis中，固然需要先发送一次请求CaptchaRequest,然后用返回的验证码
	//邮箱以及密码进行注册，后续再将账号进行返回
	RegisterRequest struct {
		Email    string `json:"email" binding:"required"`
		Captcha  string `json:"captcha"`
		Password string `json:"password"`
	}
	//注册成功之后，直接让其进行登录状态
	RegisterResponse struct {
		controller.Response
		UserID int64  `json:"user_id,omitempty"`
		Token  string `json:"token,omitempty"`
	}

	CaptchaRequest struct {
		Email string `json:"email" binding:"required"`
	}

	CaptchaResponse struct {
		controller.Response
	}

	LikeRequest struct {
		FileID string `json:"file_id" binding:"required"`
	}

	LikeResponse struct {
		controller.Response
		LikeCnt    int64 `json:"like_count,omitempty"`
		LikeStatus int64 `json:"like_status,omitempty"`
	}
)

func Login(c *gin.Context) {
	var ok bool
	var userInformation *model.User
	req := new(LoginRequest)
	res := new(LoginResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	//1:判断用户是否存在
	if ok, userInformation = user.IsExistUser(req.Username); !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeUserNotExist))
		return
	}
	//2:判断用户是否密码账号正确
	if userInformation.Password != utils.MD5(req.Password) {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidPassword))
		return
	}
	//3:返回userid 以及一个Token
	token, err := myjwt.GenerateToken(userInformation.ID, userInformation.Username)

	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	res.Success()
	res.UserID, res.Token = userInformation.ID, token
	c.JSON(http.StatusOK, res)

}

func Register(c *gin.Context) {
	var ok bool
	var userInformation *model.User
	req := new(RegisterRequest)
	res := new(RegisterResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	//1:先判断用户是否已经存在了
	if ok, _ := user.IsExistUser(req.Email); ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeUserExist))
		return
	}

	if ok, userInformation = user.Register(req.Email, req.Password, req.Captcha); !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	// 生成Token
	token, err := myjwt.GenerateToken(userInformation.ID, userInformation.Username)

	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	res.Success()
	res.UserID, res.Token = userInformation.ID, token
	c.JSON(http.StatusOK, res)
}

func HandleCaptcha(c *gin.Context) {
	req := new(CaptchaRequest)
	res := new(CaptchaResponse)
	//解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	//给service层进行处理
	if ok := user.SendCaptcha(req.Email); !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}
	//匿名字段，其实本身res.Success()调用就是res.Response.Success()
	//res.Response.Success()
	res.Success()
	c.JSON(http.StatusOK, res)
}

// 登录之后，进行点赞操作
func Like(c *gin.Context) {
	req := new(LikeRequest)
	res := new(LikeResponse)
	//解析参数
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	userID := c.GetInt64("user_id") // 从中间件 Set() 中获取
	//开始给service层进行点赞业务处理
	LikeCnt, LikeStatus, ok := user.HandleLike(userID, req.FileID)
	if !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}
	res.Success()
	res.LikeStatus, res.LikeCnt = LikeStatus, LikeCnt
	c.JSON(http.StatusOK, res)
}
