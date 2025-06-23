package music

import (
	"MyChat/controller"
	"github.com/gin-gonic/gin"
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
		FileID string `form:"uuid" binding:"required"`
	}

	MusicDownloadResponse struct {
		controller.Response
	}
)

func MusicUpload(c *gin.Context) {
	req := new(MusicUploadRequest)
	res := new(MusicUploadResponse)

}

// 下载需要传音乐的uuid，（jwt，这个在user_id中）
// 然后返回FileSize,文件，Duration
func MusicDownload(c *gin.Context) {
	req := new(MusicDownloadRequest)
	res := new(MusicDownloadResponse)
	//1:获取user_id,从/static/user_id中找是否存在file_id为前缀的文件

	//2:如果有，获取，如果没有，那么返回false，文件不存在

}
