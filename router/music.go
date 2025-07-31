package router

import (
	"MyChat/controller/music"
	"github.com/gin-gonic/gin"
)

func MusicRouter(r *gin.RouterGroup) {
	{
		r.POST("/music_upload", music.MusicUpload)
		r.GET("/music_download", music.MusicDownload)
		r.GET("/stream", music.MusicStart)
		r.GET("/rankings", music.Rankings)
		//用于”发现音乐“，懒加载音乐信息，
		r.GET("music_infos", music.MusicInfo)
	}
}
