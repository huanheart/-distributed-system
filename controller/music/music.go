package music

import (
	"MyChat/common/code"
	"MyChat/config"
	"MyChat/controller"
	"MyChat/model"
	"MyChat/service/music"
	"MyChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
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
	// GET请求需绑定form
	MusicDownloadRequest struct {
		FileID string `form:"file_id" binding:"required"`
	}

	MusicDownloadResponse struct {
		controller.Response
		FilePath string `json:"file_path" binding:"required"`
	}
	MusicFileInfoRequest struct {
		Id  int64 `json:"id" binding:"required"`
		Cnt int64 `json:"cnt" binding:"required"`
	}
	//用于初始化获取现有文件信息的响应类
	MusicFileInfoResponse struct {
		controller.Response
		MusicInfoList []controller.MusicInfo `json:"music_info_list"`
	}

	//播放音乐的请求类
	MusicStartRequest struct {
		FileID string `form:"file_id" binding:"required"`
	}
	MusicStartResponse struct {
		controller.Response
	}
	//关于音乐排行榜的响应类，不需要请求类
	//1：音乐文件图片路径  2：点赞数量

	MusicRankingsResponse struct {
		controller.Response
		MusicList []controller.MusicDetail `json:"music_list"` // 返回的音乐信息列表
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
	if err := c.ShouldBindQuery(req); err != nil {
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

	res.FilePath = utils.GetHttpPath(musicfile.FilePath)

	log.Println("res.FilePath is " + res.FilePath)
	c.JSON(http.StatusOK, res)
}

func MusicInfo(c *gin.Context) {
	var file_paths []controller.MusicInfo
	var ok bool
	req := new(MusicFileInfoRequest)
	res := new(MusicFileInfoResponse)
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}
	//获取数据库中大于当前Id的前Cnt的行，并将对应相关数据进行返回
	if file_paths, ok = music.GetMusicFilesAfterID(req.Id, req.Cnt); !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}
	//遍历当前这个file_paths中的FilePath字段，更新为http路径
	for i := range file_paths {
		file_paths[i].FilePath = utils.GetHttpPath(file_paths[i].FilePath)
	}

	res.MusicInfoList = file_paths
	//成功标志
	res.Success()
	c.JSON(http.StatusOK, res)
}

// 获取音乐排行榜数据，Get请求返回对应的
func Rankings(c *gin.Context) {
	var file_paths []controller.MusicDetail
	var ok bool
	res := new(MusicRankingsResponse)
	//获取点赞数前五的音乐 的图片文件路径，点赞数 ,默认获取RedisRankingsNum个
	if file_paths, ok = music.GetTopInformation(config.DefaultRedisKeyConfig.RedisRankingsNum); !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	res.MusicList = file_paths
	res.Success()
	c.JSON(http.StatusOK, res)
}

func MusicStart(c *gin.Context) {

	req := new(MusicStartRequest)

	//这个res是用于错误请求进行返回
	res := new(MusicStartResponse)

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	userID := c.GetInt64("user_id") // 来自 JWT 中间件

	// 1. 查询数据库获取音乐文件信息
	musicFile, ok := music.IsExistMusicFile(userID, req.FileID)
	if !ok {
		c.JSON(http.StatusOK, res.CodeOf(code.FileNotFind))
		return
	}

	filePath := musicFile.FilePath
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusOK, res.CodeOf(code.FileCannotOpen))
		return
	}
	defer file.Close()

	fi, _ := file.Stat()
	size := fi.Size()

	// 2. 解析 Range 请求
	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		// 如果不支持 Range，直接返回整段
		c.Header("Content-Type", "audio/mpeg")
		c.Header("Content-Length", fmt.Sprintf("%d", size))
		c.Header("Accept-Ranges", "bytes")
		http.ServeContent(c.Writer, c.Request, filePath, fi.ModTime(), file)
		return
	}

	// Range 格式：bytes=start-end
	var start, end int64
	_, err = fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if err != nil {
		// 若未给出 end，默认读到结尾
		_, err = fmt.Sscanf(rangeHeader, "bytes=%d-", &start)
		if err != nil {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
			return
		}
		end = size - 1
	}

	// 校验范围合法
	if start > end || end >= size {
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return
	}

	// 3. 设置 Range 响应头并返回部分内容
	length := end - start + 1

	c.Status(http.StatusPartialContent) // 206
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Content-Length", fmt.Sprintf("%d", length))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))

	// 定位到起始字节
	file.Seek(start, io.SeekStart)
	io.CopyN(c.Writer, file, length)
}
