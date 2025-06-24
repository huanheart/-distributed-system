package music

import (
	"MyChat/common/code"
	"MyChat/config"
	"MyChat/controller"
	"MyChat/model"
	"MyChat/service/music"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	//MusicUploadRequest struct {
	//}
	MusicUploadResponse struct {
		FileID string `json:"file_id" binding:"required"`
		controller.Response
	}

	// 下载需要传音乐的uuid，（jwt，这个在user_id中）
	// 然后返回FileSize,文件，Duration
	MusicDownloadRequest struct {
		FileID string `json:"file_id" binding:"required"`
	}

	MusicDownloadResponse struct {
		controller.Response
		FilePath string `json:"file_path" binding:"required"`
	}
)

//todo:后续还要提供获取所有音乐文件的接口，在用户登录界面的时候得进行一个加载

func MusicUpload(c *gin.Context) {
	//req := new(MusicUploadRequest)
	res := new(MusicUploadResponse)
	file, err := c.FormFile("file")
	if err != nil {
		res.CodeOf(code.CodeInvalidParams)
		c.JSON(http.StatusOK, res)
	}

	userID := c.GetInt64("user_id") // 从中间件 Set() 中获取
	music_file, ok := music.MusicUpload(userID, file)
	if !ok {
		res.CodeOf(code.CodeServerBusy)
		c.JSON(http.StatusOK, res)
	}
	// 8. 返回成功响应
	res.Success()
	res.FileID = music_file.UUID
	c.JSON(http.StatusOK, res)
}

// 下载需要传音乐的uuid，（jwt，这个在user_id中）
// 然后返回FileSize,文件，Duration
func MusicDownload(c *gin.Context) {
	var musicfile *model.MusicFile
	var ok bool
	req := new(MusicDownloadRequest)
	res := new(MusicDownloadResponse)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	userID := c.GetInt64("user_id") // 从中间件 Set() 中获取

	//1:从数据库中user_id,从/static/user_id中找是否存在file_id(uuid)为前缀的文件
	musicfile, ok = music.IsExistMusicFile(userID, req.FileID)
	if !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.FileNotFind))
		return
	}
	res.Success()
	res.FilePath = config.GetConfig().MusicFileIp + musicfile.FilePath
	c.JSON(http.StatusOK, res)
}
