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
	MusicUploadRequest struct {
	}
	MusicUploadResponse struct {
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

func MusicUpload(c *gin.Context) {
	req := new(MusicUploadRequest)
	res := new(MusicUploadResponse)

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
	//2:如果有，获取，如果没有，那么返回false，文件不存在

	res.Success()
	//todo:硬编码待删除
	res.FilePath = config.GetConfig().MusicFileIp + musicfile.FilePath
	c.JSON(http.StatusOK, res)
}
